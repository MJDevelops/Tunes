#pragma once

#include <vector>
#include <string>
#include <taglib/fileref.h>
#include <taglib/tpropertymap.h>

class TagFile
{
public:
    TagFile(const std::string &path);
    std::string getTitle();
    std::string getAlbum();
    std::vector<std::string> getArtists();

private:
    TagLib::FileRef f;
};