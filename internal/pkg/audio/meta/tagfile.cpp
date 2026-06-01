#include "tagfile.hpp"

#include <taglib/fileref.h>
#include <taglib/tpropertymap.h>
#include <taglib/tstringlist.h>

#include <cstdint>
#include <string>
#include <utility>
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

const std::vector<std::string>& TagFile::getAlbumArtists() {
    if (albumArtists.empty()) {
        PropertyMap props = properties();
        StringList aa = props["ALBUMARTIST"];
        for (auto& a : aa)
            this->albumArtists.push_back(a.to8Bit(true));
    }
    return albumArtists;
}

/**
 * Returns the artwork and the mime type of the artwork.
 * @return A pair with the artwork data and the mime type.
 */
std::pair<std::vector<char>, std::string> TagFile::getArtwork() {
    std::pair<std::vector<char>, std::string> vals;

    auto pics = complexProperties("PICTURE");
    if (pics.isEmpty()) return vals;

    VariantMap& map = pics.front();
    ByteVector picData = map["data"].value<ByteVector>();

    std::copy(picData.begin(), picData.end(), std::back_inserter(vals.first));
    vals.second = map["mimeType"].value<String>().to8Bit(true);

    return vals;
}