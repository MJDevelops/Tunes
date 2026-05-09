#include <stdlib.h>
#include <libavutil/avutil.h>
#include "sample_buffer.h"

SampleBuffer *sb_alloc()
{
    SampleBuffer *buf = calloc(1, sizeof(*buf));
    buf->data = av_mallocz(2 * sizeof(*buf->data));
    return buf;
}

/**
 * Frees the buffer and sets the underlying pointer to NULL.
 * @param buf The buffer to free.
 */
void sb_free(SampleBuffer **buf)
{
    if (buf && *buf)
    {
        for (int i = 0; i < 2; i++)
        {
            av_freep(&(*buf)->data[i]);
        }
        av_freep(&(*buf)->data);
        free(*buf);
        *buf = NULL;
    }
}

int16_t **sb_flush(SampleBuffer *buf)
{
    int16_t **tmp = av_mallocz(2 * sizeof(*buf));
    for (int i = 0; i < 2; i++)
    {
        tmp[i] = av_malloc(buf->channel_size * sizeof(*tmp[i]));
        memcpy(tmp[i], buf->data[i], buf->channel_size * sizeof(*buf->data[i]));
        av_freep(&buf->data[i]);
    }
    buf->channel_size = 0;
    return tmp;
}