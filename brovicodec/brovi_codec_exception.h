#ifndef _BROVI_CODEC_EXCEPTION
#define _BROVI_CODEC_EXCEPTION

#include <exception>

class BroviCodecInitException : std::exception
{
    const char *what() const throw() override
    {
        return "failed to init a brovi codec";
    }
};

class BroviCodecEncodeException : std::exception
{
    const char *what() const throw() override
    {
        return "failed to encode";
    }
};

class BroviCodecZeroFrameSizeException : std::exception
{
    const char *what() const throw() override
    {
        return "encoded frame has zero size";
    }
};

#endif