// src/utils/request.ts
import axios from 'axios';
import type {
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosRequestHeaders
} from 'axios';
import { useUserStoreWithOut } from '@/store/user';
import router from '@/router';

// 定义请求配置接口
interface RequestConfig extends AxiosRequestConfig {
  /** 是否为 FormData 类型请求 */
  isFormData?: boolean;
  /** 是否跳过认证 */
  skipAuth?: boolean;
  /** 自定义方法别名，与 Axios 的 method 一致 */
  method?: string;
}

// Token 获取器类型
type TokenGetter = () => string | undefined;

// 自定义请求错误类
export class RequestError extends Error {
  public code: number;

  constructor(
    code: number,
    message: string
  ) {
    super(message);
    this.code = code;
    this.name = 'RequestError';
  }
}

// 创建 axios 实例
const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
axiosInstance.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const originalConfig = config as InternalAxiosRequestConfig & RequestConfig;

    console.debug('请求拦截器 - 原始配置:', {
      url: originalConfig.url,
      isFormData: originalConfig.isFormData,
      dataType: originalConfig.data?.constructor?.name,
      dataIsFormData: originalConfig.data instanceof FormData,
      headers: originalConfig.headers
    });

    // 处理 URL 和 params
    if (originalConfig.params && Object.keys(originalConfig.params).length > 0) {
      const url = new URL(originalConfig.url || '', window.location.origin);
      Object.keys(originalConfig.params).forEach(key => {
        const val = originalConfig.params![key];
        if (val !== undefined && val !== null) {
          url.searchParams.append(key, String(val));
        }
      });
      config.url = url.toString().replace(window.location.origin, '');
    }

    // 处理认证令牌
    if (!originalConfig.skipAuth) {
      try {
        const userStore = useUserStoreWithOut();
        const token = userStore.getToken;
        if (token) {
          if (!config.headers) {
            config.headers = {} as AxiosRequestHeaders;
          }
          (config.headers as AxiosRequestHeaders).Authorization = `Bearer ${token}`;
        }
      } catch {
        // 忽略获取 token 失败
      }
    }

    // 处理 FormData - 增强判断逻辑
    const isFormData = originalConfig.isFormData || originalConfig.data instanceof FormData;

    if (isFormData) {
      console.debug('检测到 FormData 请求，将删除 Content-Type 头');

      // 确保 data 是 FormData
      if (originalConfig.data && !(originalConfig.data instanceof FormData)) {
        const formData = new FormData();
        Object.keys(originalConfig.data).forEach(key => {
          formData.append(key, originalConfig.data[key]);
        });
        config.data = formData;
      }

      // 删除 Content-Type 头
      if (config.headers) {
        delete (config.headers as Record<string, any>)['Content-Type'];
      }
    } else {
      console.debug('非 FormData 请求，Content-Type 保持不变');
    }

    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data: result, status } = response;

    // HTTP 状态码错误
    if (status < 200 || status >= 300) {
      if (status === 401) {
        handleAuthError();
        throw new RequestError(401, '登录已过期，请重新登录');
      }
      throw new RequestError(status, `HTTP Error: ${status}`);
    }

    // 处理业务码逻辑
    if (result.code !== undefined) {
      if (result.code !== 0) {
        if (result.code === 401) {
          handleAuthError();
          throw new RequestError(401, result.message || '登录已过期');
        }
        throw new RequestError(result.code, result.message || '请求失败');
      }
    }

    // 智能返回：有 data 字段返回 data，否则返回整个响应
    response.data = result.data !== undefined ? result.data : result;
    return response;
  },
  (error: AxiosError) => {
    if (axios.isCancel(error)) {
      throw new RequestError(504, '请求超时');
    }

    if (error.response) {
      const { status, data } = error.response;
      const responseData = data as any;

      if (status === 401) {
        handleAuthError();
        throw new RequestError(401, responseData?.message || '登录已过期');
      }

      if (responseData?.code !== undefined && responseData.code !== 0) {
        throw new RequestError(responseData.code, responseData.message || '请求失败');
      }

      throw new RequestError(status, responseData?.message || `HTTP Error: ${status}`);
    } else if (error.request) {
      throw new RequestError(500, '网络错误，请检查网络连接');
    } else {
      throw new RequestError(500, error.message || '请求配置错误');
    }
  }
);

const handleAuthError = () => {
  try {
    const userStore = useUserStoreWithOut();
    userStore.clearUser();
    router.replace('/login');
  } catch (e) {
    console.warn('Failed to clear user store');
  }
};

// 创建请求
const createRequest = async <T = any>(
  config: RequestConfig,
  getToken?: TokenGetter
): Promise<T> => {
  console.debug('createRequest 接收配置:', {
    url: config.url,
    isFormData: config.isFormData,
    dataType: config.data?.constructor?.name,
    responseType: config.responseType
  });

  const axiosConfig: RequestConfig = {
    url: config.url,
    method: config.method?.toUpperCase() || 'GET',
    headers: config.headers,
    params: config.params,
    data: config.data,
    timeout: config.timeout,
    isFormData: config.isFormData,
    skipAuth: config.skipAuth,
    responseType: config.responseType || 'json',
  };

  // 处理 FormData
  if (config.isFormData) {
    console.debug('createRequest 中检测到 isFormData=true');
    if (config.data instanceof FormData) {
      axiosConfig.data = config.data;
    } else if (config.data && typeof config.data === 'object') {
      const formData = new FormData();
      Object.keys(config.data).forEach(key => {
        formData.append(key, config.data[key]);
      });
      axiosConfig.data = formData;
    }
  }

  // 如果传入了自定义的 token 获取器
  if (getToken) {
    const token = getToken();
    if (token && !config.skipAuth) {
      axiosConfig.headers = {
        ...(axiosConfig.headers || {}),
        Authorization: `Bearer ${token}`,
      };
    }
  }

  console.debug('最终 axios 配置:', {
    url: axiosConfig.url,
    headers: axiosConfig.headers,
    dataType: axiosConfig.data?.constructor?.name
  });

  try {
    const response = await axiosInstance.request<T>(axiosConfig);
    return response.data;
  } catch (error: unknown) {
    if (error instanceof RequestError) {
      throw error;
    }
    throw new RequestError(500, error instanceof Error ? error.message : '未知错误');
  }
};

// 保持原有 API 接口
export const request = {
  get<T = any>(url: string, params?: Record<string, any>, options?: Partial<RequestConfig>) {
    return createRequest<T>({ url, params, method: 'GET', ...options });
  },
  post<T = any>(url: string, data?: any, options?: Partial<RequestConfig>) {
    return createRequest<T>({ url, data, method: 'POST', ...options });
  },
  put<T = any>(url: string, data?: any, options?: Partial<RequestConfig>) {
    return createRequest<T>({ url, data, method: 'PUT', ...options });
  },
  delete<T = any>(url: string, params?: Record<string, any>, options?: Partial<RequestConfig>) {
    return createRequest<T>({ url, params, method: 'DELETE', ...options });
  },
  patch<T = any>(url: string, data?: any, options?: Partial<RequestConfig>) {
    return createRequest<T>({ url, data, method: 'PATCH', ...options });
  },
  upload<T = any>(url: string, formData: FormData, options?: Partial<RequestConfig>) {
    return createRequest<T>({
      url,
      data: formData,
      method: 'POST',
      isFormData: true,
      ...options
    });
  },
  // 额外提供 axios 实例，方便需要更精细控制的场景
  get axiosInstance() {
    return axiosInstance;
  }
};