#include "tagfile.hpp"

#include <taglib/fileref.h>
#include <taglib/tpropertymap.h>
#include <taglib/tstringlist.h>

#include <string>
#include <vector>

using namespace TagLib;

TagFile::TagFile(const std::string& path) : FileRef(path.c_str()) {}

TagFile::~TagFile() {}

std::string TagFile::getTitle() {
    return tag()->title().to8Bit(true);
}

std::string TagFile::getAlbum() {
    return tag()->album().to8Bit(true);
}

const std::vector<std::string>& TagFile::getArtists() {
    if (artists.empty()) {
        PropertyMap props = properties();
        StringList artists = props["ARTIST"];
        for (auto& a : artists)
            this->artists.push_back(a.to8Bit(true));
    }
    return artists;
}

std::string TagFile::getArtwork() {
    return "";
}