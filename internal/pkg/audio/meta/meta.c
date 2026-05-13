#include <stdlib.h>
#include <taglib/tag_c.h>
#include "meta.h"

TagFile *tf_alloc(const char *file)
{
    TagFile *tf = malloc(sizeof(*tf));
    tf->f = taglib_file_new(file);
    tf->tag = taglib_file_tag(tf->f);
    return tf;
}

void tf_free(TagFile **tf)
{
    taglib_tag_free_strings();
    taglib_file_free((*tf)->f);
    *tf = NULL;
}

int tf_get_artists(TagFile *tf, char **arr)
{
    return -1;
}

const char *tf_get_album(TagFile *tf)
{
    return taglib_tag_album(tf->tag);
}