#pragma once

#include <libavformat/avformat.h>
#include "sample_buffer.h"

typedef struct
{
    AVCodecContext *ctx;
    AVFormatContext *fmt_ctx;
    AVPacket *pkt;
    AVFrame *frame;
    AVFrame *resampled_frame;
    const char *filename;
    const AVCodec *codec;
    int audio_stream_index;
} Decoder;

Decoder *decoder_alloc(const char *file);
void decoder_free(Decoder **);
int decode(Decoder *, SampleBuffer *, int frames);
int resample_frame_s16_planar_stereo(AVFrame *, AVFrame *);