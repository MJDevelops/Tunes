%module meta

%include "std_vector.i"
%include "std_string.i"

namespace std {
    %template(StringVector) vector<string>;
};

%{
#include "tagfile.hpp"
%}

%include "tagfile.hpp"