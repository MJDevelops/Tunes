#include "libavformat/avformat.h"
#include "libavcodec/avcodec.h"
#include "libavutil/avutil.h"
#include "decode.h"
#include <stdint.h>
#include <stdlib.h>

static int16_t *decode(const char *filename)
{
    int16_t *interleave_buf = NULL;
    int16_t *buf = NULL;
    AVCodec *codec = NULL;
    AVCodecContext *ctx = NULL;
    AVPacket *pkt = av_packet_alloc();
    AVFrame *frame = av_frame_alloc();
    AVFormatContext *fmt_ctx = avformat_alloc_context();
    int64_t size_buf = 0;
    int size_ibuf = 0;
    int audio_stream_index;
    int ret;

    avformat_open_input(&fmt_ctx, filename, NULL, NULL);
    avformat_find_stream_info(fmt_ctx, NULL);

    for (int i = 0; i < fmt_ctx->nb_streams; i++)
    {
        if (fmt_ctx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_AUDIO)
        {
            audio_stream_index = i;

            codec = avcodec_find_decoder(fmt_ctx->streams[i]->codecpar->codec_id);
            if (!codec)
            {
                fprintf(stderr, "Codec not found\n");
                return NULL;
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
            av_packet_unref(pkt);
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
                fprintf(stderr, "Unsupported audio format\n");
                goto free;
            }

            interleave_buf = av_malloc_array(frame->nb_samples * 2, sizeof(*interleave_buf));
            if (!interleave_buf)
            {
                fprintf(stderr, "Could not allocate buffer\n");
                goto free;
            }

            size_ibuf = frame->nb_samples * 2 * sizeof(*interleave_buf);

            for (int i = 0; i < frame->nb_samples; i++)
            {
                interleave_buf[2 * i] = ((int16_t *)frame->data[0])[i];
                interleave_buf[2 * i + 1] = ((int16_t *)frame->data[1])[i];
            }

            buf = av_realloc_array(buf, size_buf + frame->nb_samples * 2, sizeof(*buf));
            if (!buf)
            {
                av_freep(&interleave_buf);
                goto free;
            }

            memcpy(buf + size_buf, interleave_buf, size_ibuf);

            size_buf += frame->nb_samples * 2;

            av_freep(&interleave_buf);
            av_frame_unref(frame);
        }
        av_packet_unref(pkt);
    }
free:
    av_freep(&pkt);
    av_freep(&frame);
    avformat_close_input(&fmt_ctx);
    avcodec_free_context(&ctx);
    return buf;
}