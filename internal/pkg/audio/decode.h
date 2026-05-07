#ifndef DECODE_H
#define DECODE_H
#include <stdint.h>
#include "libavformat/avformat.h"

int64_t decode(double_t **, const char *);
int resample_frame_double_planar_stereo(AVFrame *, AVFrame *);

#endif