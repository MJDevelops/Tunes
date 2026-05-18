%module meta

%include "std_vector.i"
%include "std_string.i"
%include "std_pair.i"

namespace std {
    %template(StringVector) vector<string>;
    %template(ByteVector) vector<char>;
    %template(ByteStringPair) pair<vector<char>, string>;
};

%{
#include "tagfile.hpp"
%}

%include "tagfile.hpp"