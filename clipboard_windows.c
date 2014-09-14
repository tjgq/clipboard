#include <stdlib.h>
#include <windows.h>

char *get() {
  char *ret = NULL;
  int sz;
  HANDLE h;

  if (!OpenClipboard(NULL))
    goto done;

  h = GetClipboardData(CF_UNICODETEXT);
  if (!h)
    goto close;

  if (!GlobalLock(h))
    goto close;

  sz = WideCharToMultiByte(CP_UTF8, 0, h, -1, NULL, 0, NULL, NULL);
  if (!sz)
    goto unlock;

  ret = malloc(sz);
  if (!ret)
    goto unlock;

  sz = WideCharToMultiByte(CP_UTF8, 0, h, -1, ret, sz, NULL, NULL);
  if (!sz) {
    free(ret);
    ret = NULL;
  }

unlock:
  GlobalUnlock(h);
close:
  CloseClipboard();
done:
  return ret;
}

int set(const char *s) {
  int ret = 0;
  wchar_t *h;
  int sz;
  
  sz = MultiByteToWideChar(CP_UTF8, 0, s, -1, NULL, 0);
  if (!sz)
    goto done;
  
  h = GlobalAlloc(0, 2*sz);
  if (!h)
    goto done;
    
  sz = MultiByteToWideChar(CP_UTF8, 0, s, -1, h, sz);
  if (!sz) {
    goto dealloc;
  }
  
  if (!OpenClipboard(NULL))
    goto dealloc;
    
  if (!EmptyClipboard())
    goto dealloc;
    
  if (SetClipboardData(CF_UNICODETEXT, h)) {
    ret = 1;
    goto close;
  }
   
dealloc:
    GlobalFree(h);
close:
    CloseClipboard();
done:
  return ret;
}
