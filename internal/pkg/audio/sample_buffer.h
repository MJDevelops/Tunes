#pragma once

#include <stdint.h>

typedef struct
{
    int sample_rate;
    int64_t channel_size;
    int16_t **data;
} SampleBuffer;

SampleBuffer *sb_alloc();
void sb_free(SampleBuffer **);