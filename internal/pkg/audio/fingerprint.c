#include "fingerprint.h"

#include <chromaprint.h>
#include <stdlib.h>

#include "decoder.h"

char* fingerprint_file(const char* path) {
    char* fingerprint;
    Decoder* dec = decoder_alloc(path);
    SampleBuffer* buf = sb_alloc();
    int ret = decode(dec, buf, -1);
    if (ret < 0) {
        goto free;
    }

    int16_t* interleaved_buf = sb_interleave(buf);

    ChromaprintContext* ctx = chromaprint_new(CHROMAPRINT_ALGORITHM_DEFAULT);
    ret = chromaprint_start(ctx, dec->sample_rate, 2);
    if (ret == 0) {
        goto free;
    }

    ret = chromaprint_feed(ctx, interleaved_buf, buf->channel_size * 2);
    if (ret == 0) {
        goto free;
    }

    ret = chromaprint_finish(ctx);
    if (ret == 0) {
        goto free;
    }

    ret = chromaprint_get_fingerprint(ctx, &fingerprint);
    if (ret == 0) {
        goto free;
    }

free:
    sb_free(&buf);
    decoder_free(&dec);

    if (ctx != NULL) {
        chromaprint_free(ctx);
    }

    if (interleaved_buf) {
        free(interleaved_buf);
    }

    return fingerprint;
}