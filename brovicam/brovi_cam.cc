#include <fcntl.h>
#include <cstring>
#include <unistd.h>
#include <linux/videodev2.h>
#include <sys/mman.h>
#include "brovi_cam.h"
#include "brovi_exception.h"
#include "xioctl.h"

#define VIDEO_BUFFER_NUMBER 4

class KBroviCam
{
  public:
    explicit KBroviCam(BroviCamConfig *);
    ~KBroviCam() noexcept;
    void Close();
    void Start();
    void Stop();
    VideoBufferStatus NextBufferA();
    void NextBufferB(VideoBufferStatus);

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
    fmt.fmt.pix.width = width;
    fmt.fmt.pix.height = height;
    fmt.fmt.pix.pixelformat = V4L2_PIX_FMT_YUYV;
    fmt.fmt.pix.field = V4L2_FIELD_INTERLACED;
    if (xioctl(fd, VIDIOC_S_FMT, &fmt) == -1)
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
    if (xioctl(fd, VIDIOC_REQBUFS, &req) < 0)
    {
        throw BroviCamOpenException();
    }
}

void KBroviCam::QueryBuffers()
{
    v4l2_buffer buf;
    for (int i = 0; i < VIDEO_BUFFER_NUMBER; i++)
    {
        memset(&buf, 0, sizeof(buf));
        buf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        buf.memory = V4L2_MEMORY_MMAP;
        buf.index = i;
        if (xioctl(fd, VIDIOC_QUERYBUF, &buf) < 0)
        {
            throw BroviCamOpenException();
        }
        buffers[i].length = buf.length;
        buffers[i].start = mmap(nullptr, buf.length, PROT_READ | PROT_WRITE, MAP_SHARED, fd, buf.m.offset);
        if (buffers[i].start == MAP_FAILED)
        {
            throw BroviCamOpenException();
        }
        if (xioctl(fd, VIDIOC_QBUF, &buf) < 0)
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

KBroviCam::~KBroviCam() noexcept
{
    for (int i = 0; i < VIDEO_BUFFER_NUMBER; i++)
    {
        munmap(buffers[i].start, buffers[i].length);
    }
    delete[] buffers;
}

void KBroviCam::Close()
{
    close(fd);
}

void KBroviCam::Start()
{
    v4l2_buf_type type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    if (xioctl(fd, VIDIOC_STREAMON, &type) < 0)
    {
        throw BroviCamStartException();
    }
}

void KBroviCam::Stop()
{
    v4l2_buf_type type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    if (xioctl(fd, VIDIOC_STREAMOFF, &type) < 0)
    {
        throw BroviCamStopException();
    }
}

VideoBufferStatus KBroviCam::NextBufferA()
{
    v4l2_buffer buf;
    memset(&buf, 0, sizeof(buf));
    buf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    buf.memory = V4L2_MEMORY_MMAP;
    if (xioctl(fd, VIDIOC_DQBUF, &buf) < 0)
    {
        throw BroviCamNextBufferException();
    }
    return VideoBufferStatus{buffer : &buffers[buf.index], v4l_buf : buf};
}

void KBroviCam::NextBufferB(VideoBufferStatus status)
{
    if (xioctl(fd, VIDIOC_QBUF, &status.v4l_buf) < 0)
    {
        throw BroviCamNextBufferException();
    }
}

//---------------------------------------------

BroviCam *BroviCam_Open(BroviCamConfig *config)
{
    try
    {
        return static_cast<BroviCam *>(new KBroviCam(config));
    }
    catch (BroviCamOpenException &e)
    {
        return nullptr;
    }
}

void BroviCam_Close(BroviCam *bc)
{
    static_cast<KBroviCam *>(bc)->Close();
}

int BroviCam_Start(BroviCam *bc)
{
    try
    {
        static_cast<KBroviCam *>(bc)->Start();
        return 0;
    }
    catch (BroviCamStartException &e)
    {
        return -1;
    }
}

int BroviCam_Stop(BroviCam *bc)
{
    try
    {
        static_cast<KBroviCam *>(bc)->Stop();
        return 0;
    }
    catch (BroviCamStopException &e)
    {
        return -1;
    }
}

VideoBufferStatus BroviCam_NextBufferA(BroviCam *bc)
{
    try
    {
        return static_cast<KBroviCam *>(bc)->NextBufferA();
    }
    catch (BroviCamNextBufferException &e)
    {
        return VideoBufferStatus{buffer : nullptr};
    }
}

int BroviCam_NextBufferB(BroviCam *bc, VideoBufferStatus status)
{
    try
    {
        static_cast<KBroviCam *>(bc)->NextBufferB(status);
        return 0;
    }
    catch (BroviCamNextBufferException &e)
    {
        return -1;
    }
}
