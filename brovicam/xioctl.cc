#include <cerrno>
#include <sys/ioctl.h>
#include "xioctl.h"

int xioctl(int fh, int request, void *arg)
{
    int r;

    do
    {
        r = ioctl(fh, request, arg);
    } while (-1 == r && EINTR == errno);

    return r;
}