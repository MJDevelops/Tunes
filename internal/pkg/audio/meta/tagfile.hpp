#pragma once

#include <taglib/fileref.h>
#include <taglib/tpropertymap.h>

#include <string>
#include <utility>
#include <vector>

class TagFile : private TagLib::FileRef {
   public:
    TagFile(const std::string& path);
    ~TagFile();
    std::string getTitle();
    std::string getAlbum();
    const std::vector<std::string>& getArtists();
    const std::vector<std::string>& getAlbumArtists();
    std::pair<std::vector<char>, std::string> getArtwork();

   private:
    std::vector<std::string> artists;
    std::vector<std::string> albumArtists;
};