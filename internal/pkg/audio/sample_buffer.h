#pragma once

#include <stdint.h>

typedef struct
{
    int64_t channel_size;
    int16_t **data;
} SampleBuffer;

SampleBuffer *sb_alloc();
void sb_free(SampleBuffer **);
int16_t **sb_flush(SampleBuffer *);