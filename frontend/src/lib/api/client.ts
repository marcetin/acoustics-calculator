const API_BASE = '/api';

export interface SceneConfig {
  room: {
    widthM: number;
    lengthM: number;
    heightM: number;
  };
  surfaces: {
    frontWall: { materialId: string };
    rearWall: { materialId: string };
    leftWall: { materialId: string };
    rightWall: { materialId: string };
    ceiling: { materialId: string };
    floor: { materialId: string };
  };
  sources: {
    left: { x: number; y: number; z: number };
    right: { x: number; y: number; z: number };
  };
  receiver: {
    x: number;
    y: number;
    z: number;
  };
  viewer: {
    showGrid: boolean;
    showLabels: boolean;
    showAxes: boolean;
  };
}

export interface DerivedScene {
  surfaces: Array<{
    id: string;
    label: string;
    materialId: string;
    colorHex: string;
  }>;
  roomBounds: {
    minX: number;
    maxX: number;
    minY: number;
    maxY: number;
    minZ: number;
    maxZ: number;
  };
  leftSpeaker: { x: number; y: number; z: number };
  rightSpeaker: { x: number; y: number; z: number };
  receiver: { x: number; y: number; z: number };
}

export interface PreviewResponse {
  config: SceneConfig;
  derived: DerivedScene;
  warnings: string[];
}

export interface SuccessResponse<T> {
  ok: true;
  data: T;
  warnings: string[];
}

export interface ErrorResponse {
  ok: false;
  error: {
    code: string;
    message: string;
    fields?: Record<string, string>;
  };
}

export async function getExample(): Promise<SceneConfig> {
  const response = await fetch(`${API_BASE}/scene/example`);
  const result: SuccessResponse<SceneConfig> = await response.json();
  return result.data;
}

export async function previewScene(config: SceneConfig): Promise<PreviewResponse> {
  const response = await fetch(`${API_BASE}/scene/preview`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  });
  
  const result: SuccessResponse<PreviewResponse> | ErrorResponse = await response.json();
  
  if (!result.ok) {
    throw new Error((result as ErrorResponse).error.message);
  }
  
  return (result as SuccessResponse<PreviewResponse>).data;
}
