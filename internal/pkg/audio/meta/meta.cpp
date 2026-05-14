#include <vector>
#include <string>
#include <taglib/fileref.h>
#include <taglib/tpropertymap.h>
#include <taglib/tstringlist.h>
#include "meta.hpp"

TagFile::TagFile(const std::string &path)
{
    f = TagLib::FileRef(path.c_str());
}

std::string TagFile::getTitle()
{
    return f.tag()->title().to8Bit(true);
}

std::string TagFile::getAlbum()
{
    return f.tag()->album().to8Bit(true);
}

std::vector<std::string> TagFile::getArtists()
{
    std::vector<std::string> artistVec{};
    TagLib::PropertyMap props = f.properties();
    TagLib::StringList artists = props["ARTIST"];

    for (auto &a : artists)
        artistVec.push_back(a.to8Bit(true));

    return artistVec;
}