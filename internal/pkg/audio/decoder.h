#pragma once

#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>

#include "sample_buffer.h"

typedef struct
{
    AVCodecContext* ctx;
    AVFormatContext* fmt_ctx;
    AVPacket* pkt;
    AVFrame* frame;
    AVFrame* resampled_frame;

    /** Filename of the audio file */
    const char* filename;

    /** Codec of the audio stream */
    const AVCodec* codec;

    /** Index of the audio stream in the format context */
    int audio_stream_index;

    /** Sample rate of the audio stream */
    int sample_rate;

    /** Duration in samples */
    int64_t duration;
} Decoder;

Decoder* decoder_alloc(const char* file);
void decoder_free(Decoder**);
int decode(Decoder*, SampleBuffer*, int frames);
int resample_frame_s16_planar_stereo(AVFrame*, AVFrame*);
int decoder_seek(Decoder*, int64_t offset);