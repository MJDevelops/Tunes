#include <stdint.h>
#include <stdlib.h>
#include "libavformat/avformat.h"
#include "libavcodec/avcodec.h"
#include "libavutil/avutil.h"
#include "libswresample/swresample.h"
#include "decode.h"

SampleBuffer *sb_alloc()
{
    SampleBuffer *buf = calloc(1, sizeof(*buf));
    buf->data = av_mallocz(2 * sizeof(*buf->data));
    return buf;
}

/** Frees the buffer and sets the underlying pointer to NULL. 
 * @param buf The buffer to free.
 * **/
void sb_free(SampleBuffer **buf)
{
    if (buf && *buf)
    {
        for (int i = 0; i < 2; i++)
        {
            av_freep((*buf)->data[i]);
        }
        av_freep((*buf)->data);
        free(*buf);
        *buf = NULL;
    }
}

/**
 * Decodes the specified audio file and returns the samples in signed 16 bit
 * format with a stereo channel layout.
 * @param buf The buffer to write the data to, must have a size of 2.
 * @param filename The file to decode.
 * @return The number of samples in each channel if successful, < 0 if a error occured or the specified array is not
 * of size 2.
 * **/
int decode(SampleBuffer *buf, const char *filename)
{
    if (!buf->data)
    {
        fprintf(stderr, "Supplied buffer is invalid.\n");
        return -1;
    }

    AVCodecContext *ctx = NULL;
    AVPacket *pkt = av_packet_alloc();
    AVFrame *frame = av_frame_alloc();
    AVFrame *resampled_frame = av_frame_alloc();
    AVFormatContext *fmt_ctx = avformat_alloc_context();
    const AVCodec *codec;
    int audio_stream_index;
    int ret;

    ret = avformat_open_input(&fmt_ctx, filename, NULL, NULL);
    if (ret < 0)
    {
        fprintf(stderr, "Couldn't open specified file.\n");
        goto free;
    }

    ret = avformat_find_stream_info(fmt_ctx, NULL);
    if (ret < 0)
    {
        fprintf(stderr, "Couldn't read stream info.\n");
        goto free;
    }

    for (int i = 0; i < fmt_ctx->nb_streams; i++)
    {
        if (fmt_ctx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_AUDIO)
        {
            audio_stream_index = i;

            codec = avcodec_find_decoder(fmt_ctx->streams[i]->codecpar->codec_id);
            if (!codec)
            {
                fprintf(stderr, "Codec not found\n");
                ret = -1;
                goto free;
            }

            ctx = avcodec_alloc_context3(codec);
            avcodec_parameters_to_context(ctx, fmt_ctx->streams[i]->codecpar);
            break;
        }
    }

    avcodec_open2(ctx, codec, NULL);

    while (1)
    {
        ret = av_read_frame(fmt_ctx, pkt);
        if (ret == AVERROR_EOF)
        {
            break;
        }

        if (pkt->stream_index != audio_stream_index)
        {
            av_packet_unref(pkt);
            continue;
        }

        ret = avcodec_send_packet(ctx, pkt);
        if (ret < 0)
        {
            fprintf(stderr, "Error sending packet: %s\n", av_err2str(ret));
            goto free;
        }

        while (ret >= 0)
        {
            ret = avcodec_receive_frame(ctx, frame);
            if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
            {
                break;
            }
            else if (ret < 0)
            {
                fprintf(stderr, "Error receiving frame: %s\n", av_err2str(ret));
                break;
            }

            if (frame->ch_layout.nb_channels != 2 || frame->format != AV_SAMPLE_FMT_S16P)
            {
                ret = resample_frame_s16_planar_stereo(resampled_frame, frame);
                if (ret < 0)
                {
                    fprintf(stderr, "Couldn't resample frame: %s\n", av_err2str(ret));
                    goto free;
                }
                av_frame_unref(frame);
                av_frame_move_ref(frame, resampled_frame);
            }

            if (buf->sample_rate == 0)
            {
                buf->sample_rate = frame->sample_rate;
            }

            for (int i = 0; i < 2; i++)
            {
                buf->data[i] = av_realloc(buf->data[i], (buf->channel_size + frame->nb_samples) * sizeof(*buf->data[i]));
                if (!buf->data[i])
                {
                    fprintf(stderr, "Couldn't allocate channel buffer.\n");
                    ret = -1;
                    goto free;
                }
            }

            memcpy(buf->data[0] + buf->channel_size, (int16_t *)frame->data[0], frame->nb_samples * sizeof(*buf->data[0]));
            memcpy(buf->data[1] + buf->channel_size, (int16_t *)frame->data[1], frame->nb_samples * sizeof(*buf->data[1]));

            buf->channel_size += frame->nb_samples;

            av_frame_unref(frame);
        }
        av_packet_unref(pkt);
    }
free:
    av_freep(&pkt);
    av_freep(&frame);
    av_freep(&resampled_frame);
    avformat_close_input(&fmt_ctx);
    
    if (ctx)
    {
        avcodec_free_context(&ctx);
    }

    return ret;
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