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

export interface SplitSeries {
  key: string
  dataPoints: AggregatedDataPoint[]
}

export interface GetMeasurementDataSplitResponse {
  name: string
  splitBy: string
  series: SplitSeries[]
}

export type GetMeasurementDataResult = GetMeasurementDataResponse | GetMeasurementDataSplitResponse

export function isSplitResponse(
  response: GetMeasurementDataResult
): response is GetMeasurementDataSplitResponse {
  return 'series' in response
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
      splitBy?: string
    }
  ): Promise<GetMeasurementDataResult> {
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

    // Add splitBy parameter if specified
    if (params.splitBy) {
      queryParams.splitBy = params.splitBy
    }

    return api.get(`/products/${productId}/measurements/${encodeURIComponent(name)}/data`, {
      params: queryParams,
    })
  },
}
