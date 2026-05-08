#ifndef DECODE_H
#define DECODE_H
#include <stdint.h>
#include "libavformat/avformat.h"

int64_t decode(double_t **, const char *);
int resample_frame_double_planar_stereo(AVFrame *, AVFrame *);
void free_sample_buffer(void **, int, int64_t);

#endif