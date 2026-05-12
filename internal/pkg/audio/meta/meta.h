#pragma once

#include <string>
#include <taglib/tag.h>
#include <taglib/fileref.h>

class TagFile
{
public:
    void setAlbum(std::string name);
    std::string getAlbum();

private:
    TagLib::FileRef f;
};