#ifndef _EXCEPTION_H
#define _EXCEPTION_H

#include <exception>

class BroviCamOpenException : public std::exception
{
    const char *what() const throw() override
    {
        return "Cannot open a BroviCam instance";
    }
};

class BroviCamStartException : public std::exception
{
    const char *what() const throw() override
    {
        return "Cannot start video stream";
    }
};

class BroviCamStopException : std::exception
{
    const char *what() const throw() override
    {
        return "Cannot stop video stream";
    }
};

#endif
