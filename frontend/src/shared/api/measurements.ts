import { api } from './client'

export interface MeasurementSummary {
  name: string
  metadataKeys: string[]
}

export interface MetadataValues {
  key: string
  values: string[]
}

export interface AggregatedDataPoint {
  date: string
  sum: number
  count: number
}

export interface ListMeasurementNamesResponse {
  measurements: MeasurementSummary[]
}

export interface GetMetadataValuesResponse {
  metadata: MetadataValues[]
}

export interface GetMeasurementDataResponse {
  name: string
  dataPoints: AggregatedDataPoint[]
}

export const measurementsApi = {
  listNames(productId: string): Promise<ListMeasurementNamesResponse> {
    return api.get(`/products/${productId}/measurements`)
  },

  getMetadataValues(productId: string, name: string): Promise<GetMetadataValuesResponse> {
    return api.get(`/products/${productId}/measurements/${encodeURIComponent(name)}/metadata`)
  },

  getData(
    productId: string,
    name: string,
    params: {
      start: string
      end: string
      metadata?: Record<string, string>
    }
  ): Promise<GetMeasurementDataResponse> {
    const queryParams: Record<string, string> = {
      start: params.start,
      end: params.end,
    }

    // Add metadata filters with "metadata." prefix
    if (params.metadata) {
      for (const [key, value] of Object.entries(params.metadata)) {
        queryParams[`metadata.${key}`] = value
      }
    }

    return api.get(`/products/${productId}/measurements/${encodeURIComponent(name)}/data`, {
      params: queryParams,
    })
  },
}
