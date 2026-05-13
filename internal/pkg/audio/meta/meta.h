#pragma once

#include <taglib/tag_c.h>

typedef struct
{
    TagLib_Tag *tag;
    TagLib_File *f;
} TagFile;

TagFile *tf_alloc(const char *path);
const char *tf_get_album(TagFile *);
int tf_get_artists(TagFile *, char **);
void tf_free(TagFile **);