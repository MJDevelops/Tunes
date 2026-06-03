#include "sample_buffer.h"

#include <libavutil/avutil.h>
#include <stdlib.h>

SampleBuffer* sb_alloc() {
    SampleBuffer* buf = calloc(1, sizeof(*buf));
    buf->data = av_mallocz(2 * sizeof(*buf->data));
    return buf;
}

/**
 * Frees the buffer and sets the underlying pointer to NULL.
 * @param buf The buffer to free.
 */
void sb_free(SampleBuffer** buf) {
    if (buf && *buf) {
        for (int i = 0; i < 2; i++) {
            av_freep(&(*buf)->data[i]);
        }
        av_freep(&(*buf)->data);
        free(*buf);
        *buf = NULL;
    }
}

int16_t** sb_flush(SampleBuffer* buf) {
    int16_t** tmp = av_mallocz(2 * sizeof(*tmp));
    for (int i = 0; i < 2; i++) {
        tmp[i] = av_malloc(buf->channel_size * sizeof(*tmp[i]));
        memcpy(tmp[i], buf->data[i], buf->channel_size * sizeof(*buf->data[i]));
        av_freep(&buf->data[i]);
    }
    buf->channel_size = 0;
    return tmp;
}

int16_t* sb_interleave(SampleBuffer* buf) {
    int16_t* interleaved_buf = malloc(2 * buf->channel_size * sizeof(*interleaved_buf));
    for (int i = 0; i < buf->channel_size; i++) {
        interleaved_buf[i * 2] = buf->data[0][i];
        interleaved_buf[i * 2 + 1] = buf->data[1][i];
    }
    return interleaved_buf;
}