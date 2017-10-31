#include <cstdint>
#include <x264.h>
#include "brovi_codec.h"
#include "brovi_codec_exception.h"

class KBroviCodec
{
  public:
    explicit KBroviCodec(BroviCodecConfig);
    ~KBroviCodec();
    H264Frame EncodeFrame(void *);
    int FlushDelayedFrame(H264Frame *);

  private:
    const int width = 640, height = 480; //with defaults
    const int PixelNum() const
    {
        return width * height;
    }

    x264_param_t param;
    x264_picture_t pic_in, pic_out;
    x264_t *h;
    x264_nal_t *nal;
    int i_nal;

    long pts = 0;
};

KBroviCodec::KBroviCodec(BroviCodecConfig config) : width(config.width), height(config.height)
{
    if (x264_param_default_preset(&param, "medium", nullptr) < 0)
        throw BroviCodecInitException();

    param.i_csp = X264_CSP_YV16;
    param.i_width = width;
    param.i_height = height;
    param.b_vfr_input = 0;
    param.b_repeat_headers = 1;
    param.b_annexb = 1;

    if (x264_param_apply_profile(&param, "high422") < 0)
        throw BroviCodecInitException();

    if (x264_picture_alloc(&pic_in, param.i_csp, param.i_width, param.i_height) < 0)
        throw BroviCodecInitException();

    h = x264_encoder_open(&param);
    if (!h)
        x264_picture_clean(&pic_in);
}

KBroviCodec::~KBroviCodec()
{
    x264_encoder_close(h);
    x264_picture_clean(&pic_in);
}

H264Frame KBroviCodec::EncodeFrame(void *frame)
{
    int index_y = 0, index_u = 0, index_v = 0;
    uint8_t *in = static_cast<uint8_t *>(frame);
    uint8_t *y = pic_in.img.plane[0], *u = pic_in.img.plane[1], *v = pic_in.img.plane[2];
    const int num = PixelNum() * 2;

    for (int i = 0; i < num; i += 4)
    {
        *(y + (index_y++)) = *(in + i);
        *(u + (index_u++)) = *(in + i + 1);
        *(y + (index_y++)) = *(in + i + 2);
        *(v + (index_v++)) = *(in + i + 3);
    }
    pic_in.i_pts = pts++;

    int frame_size = x264_encoder_encode(h, &nal, &i_nal, &pic_in, &pic_out);
    if (frame_size < 0)
        throw BroviCodecEncodeException();
    else if (!frame_size)
    {
        throw BroviCodecZeroFrameSizeException();
    }
    return H264Frame{data : nal->p_payload, size : frame_size};
}

int KBroviCodec::FlushDelayedFrame(H264Frame *out)
{
    int ret = x264_encoder_delayed_frames(h);
    if (ret)
    {
        out->size = x264_encoder_encode(h, &nal, &i_nal, nullptr, &pic_out);
        if (out->size < 0)
            throw BroviCodecEncodeException();
        else if (!out->size)
            throw BroviCodecZeroFrameSizeException();
        out->data = nal->p_payload;
    }
    return ret;
}

//-------------------------------------------------------------------------

BroviCodec *BroviCodec_New(BroviCodecConfig config)
{
    BroviCodec *bc;
    try
    {
        bc = static_cast<BroviCodec *>(new KBroviCodec(config));
    }
    catch (BroviCodecInitException &e)
    {
        bc = nullptr;
    }
    return bc;
}

void BroviCodec_Close(BroviCodec *bc)
{
    delete static_cast<KBroviCodec *>(bc);
}

int BroviCodec_EncodeFrame(BroviCodec *bc, void *frame, H264Frame *out)
{
    try
    {
        *out = static_cast<KBroviCodec *>(bc)->EncodeFrame(frame);
        return 0;
    }
    catch (BroviCodecEncodeException &e)
    {
        return BROVI_CODEC_ENCODE_ERR;
    }
    catch (BroviCodecZeroFrameSizeException &ze)
    {
        return BROVI_CODEC_ZERO_SIZE_ERR;
    }
}

int BroviCodec_FlushDelayedFrame(BroviCodec *bc, H264Frame *out)
{
    try
    {
        return static_cast<KBroviCodec *>(bc)->FlushDelayedFrame(out);
    }
    catch (BroviCodecEncodeException &e)
    {
        return BROVI_CODEC_ENCODE_ERR;
    }
    catch (BroviCodecZeroFrameSizeException &ze)
    {
        return BROVI_CODEC_ZERO_SIZE_ERR;
    }
}