#ifndef DECODE_H
#define DECODE_H
#include <stdint.h>
#include "libavformat/avformat.h"

typedef struct
{
    int sample_rate;
    int64_t channel_size;
    int16_t **data;
} SampleBuffer;

SampleBuffer *sb_alloc();
void sb_free(SampleBuffer *);
int decode(SampleBuffer *, const char *);
int resample_frame_s16_planar_stereo(AVFrame *, AVFrame *);

#endif