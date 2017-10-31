#ifndef _BROVI_CODEC
#define _BROVI_CODEC

#ifdef __cplusplus
extern "C" {
#endif

const int BROVI_CODEC_ENCODE_ERR = -1;
const int BROVI_CODEC_ZERO_SIZE_ERR = -2;

typedef struct
{
    int width;
    int height;
} BroviCodecConfig;

typedef struct
{
    void *data;
    int size;
} H264Frame;

typedef void BroviCodec;

BroviCodec *BroviCodec_New(BroviCodecConfig);
void BroviCodec_Close(BroviCodec *);
int BroviCodec_EncodeFrame(BroviCodec *, void *, H264Frame *);
int BroviCodec_FlushDelayedFrame(BroviCodec *, H264Frame *);

#ifdef __cplusplus
}
#endif

#endif