#include <stdint.h>
#include <stdlib.h>
#include <stdbool.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libswresample/swresample.h>
#include "decoder.h"
#include "sample_buffer.h"

Decoder *decoder_alloc(const char *filename)
{
    int ret = 0;
    Decoder *dec = calloc(1, sizeof(*dec));

    ret = avformat_open_input(dec->fmt_ctx, dec->filename, NULL, NULL);
    if (ret < 0)
    {
        fprintf(stderr, "Couldn't open specified file.\n");
        free(dec);
        return NULL;
    }

    ret = avformat_find_stream_info(dec->fmt_ctx, NULL);
    if (ret < 0)
    {
        fprintf(stderr, "Couldn't read stream info.\n");
        free(dec);
        return NULL;
    }

    for (int i = 0; i < dec->fmt_ctx->nb_streams; i++)
    {
        if (dec->fmt_ctx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_AUDIO)
        {
            dec->audio_stream_index = i;

            dec->codec = avcodec_find_decoder(dec->fmt_ctx->streams[i]->codecpar->codec_id);
            if (!dec->codec)
            {
                fprintf(stderr, "Codec not found\n");
                return NULL;
            }

            dec->ctx = avcodec_alloc_context3(dec->codec);
            avcodec_parameters_to_context(dec->ctx, dec->fmt_ctx->streams[i]->codecpar);
            break;
        }

        avcodec_open2(dec->ctx, dec->codec, NULL);
    }

    dec->frame = av_frame_alloc();
    dec->pkt = av_packet_alloc();
    dec->resampled_frame = av_frame_alloc();
    dec->filename = filename;

    return dec;
}

void decoder_free(Decoder **dec)
{
    av_freep(&(*dec)->frame);
    av_freep(&(*dec)->pkt);
    av_freep(&(*dec)->resampled_frame);

    if ((*dec)->fmt_ctx)
        avformat_close_input(&(*dec)->fmt_ctx);

    if (avcodec_is_open((*dec)->ctx))
        avcodec_free_context(&(*dec)->ctx);

    free(*dec);
    *dec = NULL;
}

// TODO: refactor this
/**
 * Decodes the specified audio file and returns the samples in signed 16 bit
 * format with a stereo channel layout.
 * @param dec The decoder which will be decoded with
 * @param buf The buffer to write the data to, must have a size of 2.
 * @param frames The number of frames to read into the sample buffer. If set to -1,
 * the file will be read until EOF
 * @return The number of frames read if successful, < 0 if an error occured.
 * **/
int decode(Decoder *dec, SampleBuffer *buf, int frames)
{
    int read = 0;
    int ret = 0;

    if (!buf->data)
    {
        fprintf(stderr, "Supplied buffer is invalid.\n");
        return -1;
    }

    while (true)
    {
        if (dec->pkt->data != NULL && dec->pkt->size > 0)
            goto receive;

        ret = av_read_frame(dec->fmt_ctx, dec->pkt);
        if (ret == AVERROR_EOF)
            break;

        if (dec->pkt->stream_index != dec->audio_stream_index)
        {
            av_packet_unref(dec->pkt);
            continue;
        }

        ret = avcodec_send_packet(dec->ctx, dec->pkt);
        if (ret < 0)
        {
            fprintf(stderr, "Error sending packet: %s\n", av_err2str(ret));
            return ret;
        }
    receive:
        while (true)
        {
            ret = avcodec_receive_frame(dec->ctx, dec->frame);
            if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
                break;
            else if (ret < 0)
            {
                fprintf(stderr, "Error receiving frame: %s\n", av_err2str(ret));
                break;
            }

            if (dec->frame->ch_layout.nb_channels != 2 || dec->frame->format != AV_SAMPLE_FMT_S16P)
            {
                ret = resample_frame_s16_planar_stereo(dec->resampled_frame, dec->frame);
                if (ret < 0)
                {
                    fprintf(stderr, "Couldn't resample frame: %s\n", av_err2str(ret));
                    av_frame_unref(dec->frame);
                    return ret;
                }
                av_frame_unref(dec->frame);
                av_frame_move_ref(dec->frame, dec->resampled_frame);
            }

            if (buf->sample_rate == 0)
                buf->sample_rate = dec->frame->sample_rate;

            for (int i = 0; i < 2; i++)
            {
                buf->data[i] = av_realloc(buf->data[i], (buf->channel_size + dec->frame->nb_samples) * sizeof(*buf->data[i]));
                if (!buf->data[i])
                {
                    fprintf(stderr, "Couldn't allocate channel buffer.\n");
                    av_frame_unref(dec->frame);
                    av_frame_unref(dec->resampled_frame);
                    return -1;
                }
            }

            memcpy(buf->data[0] + buf->channel_size, (int16_t *)dec->frame->data[0], dec->frame->nb_samples * sizeof(*buf->data[0]));
            memcpy(buf->data[1] + buf->channel_size, (int16_t *)dec->frame->data[1], dec->frame->nb_samples * sizeof(*buf->data[1]));

            buf->channel_size += dec->frame->nb_samples;

            av_frame_unref(dec->frame);

            read += 1;

            if (frames > 0 && read == frames)
            {
                av_frame_unref(dec->resampled_frame);
                goto quit;
            }
        }
        av_packet_unref(dec->pkt);
    }
quit:
    return read;
}

/**
 * Resamples an AVFrame to signed 16 bit format with a stereo channel layout.
 * @param resampled_frame The frame to resample to. Must be unreferenced.
 * @param frame The frame to resample.
 * @return >= 0 if successful, a negative AVERROR otherwise.
 **/
int resample_frame_s16_planar_stereo(AVFrame *resampled_frame, AVFrame *frame)
{
    struct SwrContext *swr_ctx = NULL;
    AVChannelLayout out_ch = AV_CHANNEL_LAYOUT_STEREO;
    int ret = 0;

    if (ret == AVERROR(EINVAL))
    {
        return ret;
    }

    ret = swr_alloc_set_opts2(
        &swr_ctx,
        &out_ch,
        AV_SAMPLE_FMT_S16P,
        frame->sample_rate,
        &frame->ch_layout,
        frame->format,
        frame->sample_rate,
        0,
        NULL);
    if (ret < 0)
    {
        return ret;
    }

    resampled_frame->sample_rate = frame->sample_rate;
    resampled_frame->ch_layout = out_ch;
    resampled_frame->format = AV_SAMPLE_FMT_S16P;

    swr_init(swr_ctx);
    swr_convert_frame(swr_ctx, resampled_frame, frame);
    swr_free(&swr_ctx);

    return 0;
}