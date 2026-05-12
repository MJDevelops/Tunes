#include <string>
#include <taglib/tag.h>
#include "meta.h"

TagFile::TagFile(std::string path)
{
    f = TagLib::FileRef(path.c_str());
}

void TagFile::setAlbum(std::string name)
{
    f.tag()->setAlbum(name);
    f.save();
}

std::string TagFile::getAlbum()
{
    return f.tag()->album().to8Bit(true);
}