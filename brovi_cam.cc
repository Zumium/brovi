#include <fcntl.h>
#include <cstring>
#include <unistd.h>
#include <sys/ioctl.h>
#include <linux/videodev2.h>
#include "brovi_cam.h"

class KBroviCam
{
  public:
    explicit KBroviCam(BroviCamConfig *);
    void Close();

  private:
    int fd;
};

KBroviCam::KBroviCam(BroviCamConfig *config)
{
    //open camera file
    fd = open(config->devfile, O_RDWR, 0);
    //set video format
    v4l2_format fmt;
    memset(&fmt, 0, sizeof(fmt));
    fmt.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    fmt.fmt.pix.width = config->width;
    fmt.fmt.pix.height = config->height;
    fmt.fmt.pix.pixelformat = V4L2_PIX_FMT_YUYV;
    fmt.fmt.pix.field = V4L2_FIELD_INTERLACED;
}

void KBroviCam::Close()
{
    close(fd);
}

//---------------------------------------------

BroviCam *BroviCam_Open(BroviCamConfig *config)
{
    return static_cast<BroviCam *>(new KBroviCam(config));
}

void Brovi_Close(BroviCam *bc)
{
    KBroviCam *kbc = static_cast<KBroviCam *>(bc);
    kbc->Close();
    delete bc;
}
