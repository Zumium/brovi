#include <fcntl.h>
#include <cstring>
#include <unistd.h>
#include <sys/ioctl.h>
#include <linux/videodev2.h>
#include <sys/mman.h>
#include "brovi_cam.h"
#include "brovi_exception.h"

#define VIDEO_BUFFER_NUMBER 4

struct VideoBuffer
{
    void *start;
    std::size_t length;
};

class KBroviCam
{
  public:
    explicit KBroviCam(BroviCamConfig *);
    ~KBroviCam();
    void Close();

  private:
    int fd;
    VideoBuffer *buffers;

    void OpenCamera(const char *);
    void SetFormat(int width, int height);
    void RequestBuffers();
    void QueryBuffers();
};

void KBroviCam::OpenCamera(const char *devfile)
{
    //open camera file
    fd = open(devfile, O_RDWR, 0);
    if (fd < 0)
    {
        throw BroviCamOpenException();
    }
}

void KBroviCam::SetFormat(int width, int height)
{
    //set video format
    v4l2_format fmt;
    memset(&fmt, 0, sizeof(fmt));
    fmt.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    fmt.fmt.pix.width = config->width;
    fmt.fmt.pix.height = config->height;
    fmt.fmt.pix.pixelformat = V4L2_PIX_FMT_YUYV;
    fmt.fmt.pix.field = V4L2_FIELD_INTERLACED;
    if (ioctl(fd, VIDIO_S_FMT, &fmt) == -1)
    {
        throw BroviCamOpenException();
    }
}

void KBroviCam::RequestBuffers()
{
    v4l2_requestbuffers req;
    req.count = VIDEO_BUFFER_NUMBER;
    req.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    req.memory = V4L2_MEMORY_MMAP;
    if (ioctl(fd, VIDIO_REQBUFS, &req) < 0)
    {
        throw BroviCamOpenException();
    }
}

KBroviCam::QueryBuffers()
{
    v4l2_buffer buf;
    for (int i = 0; i < VIDEO_BUFFER_NUMBER; i++)
    {
        memset(&buf, 0, sizeof(buf));
        buf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        buf.memory = V4L2_MEMORY_MMAP;
        buf.index = i;
        if (ioctl(fd, VIDIO_QUERYBUF, &buf) < 0)
        {
            throw BroviCamOpenException();
        }
        buffers[i].length = buf.length;
        buffers[i].start = mmap(nullptr, buf.length, PROT_READ | PROT_WRITE, MAP_SHARED, fd, buf.m.offset);
        if (buffers[i].start == MAP_FAILED)
        {
            throw BroviCamOpenException();
        }
        if (ioctl(fd, VIDIO_QBUF, &buf) < 0)
        {
            throw BroviCamOpenException();
        }
    }
}

KBroviCam::KBroviCam(BroviCamConfig *config) : buffers(new VideoBuffer[VIDEO_BUFFER_NUMBER])
{
    OpenCamera(config->devfile);
    SetFormat(config->width, config->height);
    RequestBuffers();
    QueryBuffers();
}

KBroviCam::~KBroviCam()
{
    delete[] buffers;
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
