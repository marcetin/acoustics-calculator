import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
import type { DerivedScene } from '../api/client';

export class RoomViewer {
  private scene: THREE.Scene;
  private camera: THREE.PerspectiveCamera;
  private renderer: THREE.WebGLRenderer;
  private container: HTMLElement;
  private derivedScene: DerivedScene | null = null;
  private gridHelper: THREE.GridHelper | null = null;
  private axesHelper: THREE.AxesHelper | null = null;
  private labels: THREE.Sprite[] = [];
  private controls: OrbitControls | null = null;

  constructor(container: HTMLElement) {
    this.container = container;
    
    // Scene
    this.scene = new THREE.Scene();
    this.scene.background = new THREE.Color(0x1a1a1a);

    // Camera
    this.camera = new THREE.PerspectiveCamera(
      60,
      container.clientWidth / container.clientHeight,
      0.1,
      1000
    );
    this.camera.position.set(10, 10, 8);
    this.camera.lookAt(0, 0, 0);

    // Renderer
    this.renderer = new THREE.WebGLRenderer({ antialias: true });
    this.renderer.setSize(container.clientWidth, container.clientHeight);
    this.renderer.setPixelRatio(window.devicePixelRatio);
    container.appendChild(this.renderer.domElement);

    // Lights
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.6);
    this.scene.add(ambientLight);

    const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
    directionalLight.position.set(10, 20, 10);
    this.scene.add(directionalLight);

    // Handle resize
    window.addEventListener('resize', () => this.onResize());

    // OrbitControls
    this.controls = new OrbitControls(this.camera, this.renderer.domElement);
    this.controls.enableDamping = true;
    this.controls.dampingFactor = 0.05;
    this.controls.minDistance = 1;
    this.controls.maxDistance = 100;
  }

  private onResize() {
    if (!this.container) return;
    this.camera.aspect = this.container.clientWidth / this.container.clientHeight;
    this.camera.updateProjectionMatrix();
    this.renderer.setSize(this.container.clientWidth, this.container.clientHeight);
    if (this.controls) {
      this.controls.update();
    }
  }

  updateScene(derivedScene: DerivedScene, showGrid: boolean, showLabels: boolean, showAxes: boolean) {
    // Clear existing meshes
    this.clearScene();

    this.derivedScene = derivedScene;

    // Create room shell
    this.createRoomShell(derivedScene);

    // Create speakers
    this.createSpeaker(derivedScene.leftSpeaker, 0x4CAF50);
    this.createSpeaker(derivedScene.rightSpeaker, 0x2196F3);

    // Create receiver
    this.createReceiver(derivedScene.receiver);

    // Add grid if enabled
    if (showGrid) {
      this.gridHelper = new THREE.GridHelper(
        Math.max(derivedScene.roomBounds.maxX, derivedScene.roomBounds.maxY) * 2,
        20,
        0x444444,
        0x333333
      );
      this.scene.add(this.gridHelper);
    }

    // Add axes if enabled
    if (showAxes) {
      this.axesHelper = new THREE.AxesHelper(3);
      this.scene.add(this.axesHelper);
    }

    // Add labels if enabled
    if (showLabels) {
      this.createLabels(derivedScene);
    }

    // Adjust camera to fit room
    this.fitCameraToRoom(derivedScene);
  }

  private clearScene() {
    const toRemove: THREE.Object3D[] = [];
    this.scene.traverse((child) => {
      if (child instanceof THREE.Mesh || child instanceof THREE.GridHelper || child instanceof THREE.AxesHelper) {
        toRemove.push(child);
      }
    });
    toRemove.forEach((obj) => this.scene.remove(obj));
    this.labels.forEach((label) => this.scene.remove(label));
    this.labels = [];
  }

  private createRoomShell(derivedScene: DerivedScene) {
    const { roomBounds, surfaces } = derivedScene;
    const w = roomBounds.maxX - roomBounds.minX;
    const l = roomBounds.maxY - roomBounds.minY;
    const h = roomBounds.maxZ - roomBounds.minZ;

    // Create semi-transparent walls
    const wallMaterial = (colorHex: string) =>
      new THREE.MeshLambertMaterial({
        color: colorHex,
        transparent: true,
        opacity: 0.7,
        side: THREE.DoubleSide,
      });

    // Front wall (y = 0)
    const frontWall = new THREE.Mesh(
      new THREE.PlaneGeometry(w, h),
      wallMaterial(surfaces.find((s) => s.id === 'front')?.colorHex || '#CCCCCC')
    );
    frontWall.position.set(w / 2, 0, h / 2);
    frontWall.rotation.x = Math.PI / 2;
    this.scene.add(frontWall);

    // Rear wall (y = l)
    const rearWall = new THREE.Mesh(
      new THREE.PlaneGeometry(w, h),
      wallMaterial(surfaces.find((s) => s.id === 'rear')?.colorHex || '#CCCCCC')
    );
    rearWall.position.set(w / 2, l, h / 2);
    rearWall.rotation.x = -Math.PI / 2;
    this.scene.add(rearWall);

    // Left wall (x = 0)
    const leftWall = new THREE.Mesh(
      new THREE.PlaneGeometry(l, h),
      wallMaterial(surfaces.find((s) => s.id === 'left')?.colorHex || '#CCCCCC')
    );
    leftWall.position.set(0, l / 2, h / 2);
    leftWall.rotation.y = Math.PI / 2;
    this.scene.add(leftWall);

    // Right wall (x = w)
    const rightWall = new THREE.Mesh(
      new THREE.PlaneGeometry(l, h),
      wallMaterial(surfaces.find((s) => s.id === 'right')?.colorHex || '#CCCCCC')
    );
    rightWall.position.set(w, l / 2, h / 2);
    rightWall.rotation.y = -Math.PI / 2;
    this.scene.add(rightWall);

    // Ceiling (z = h)
    const ceiling = new THREE.Mesh(
      new THREE.PlaneGeometry(w, l),
      wallMaterial(surfaces.find((s) => s.id === 'ceiling')?.colorHex || '#CCCCCC')
    );
    ceiling.position.set(w / 2, l / 2, h);
    ceiling.rotation.x = Math.PI;
    this.scene.add(ceiling);

    // Floor (z = 0)
    const floor = new THREE.Mesh(
      new THREE.PlaneGeometry(w, l),
      wallMaterial(surfaces.find((s) => s.id === 'floor')?.colorHex || '#CCCCCC')
    );
    floor.position.set(w / 2, l / 2, 0);
    this.scene.add(floor);
  }

  private createSpeaker(position: { x: number; y: number; z: number }, color: number) {
    const geometry = new THREE.BoxGeometry(0.3, 0.3, 0.4);
    const material = new THREE.MeshLambertMaterial({ color });
    const speaker = new THREE.Mesh(geometry, material);
    speaker.position.set(position.x, position.y, position.z);
    this.scene.add(speaker);
  }

  private createReceiver(position: { x: number; y: number; z: number }) {
    const geometry = new THREE.SphereGeometry(0.2, 16, 16);
    const material = new THREE.MeshLambertMaterial({ color: 0xFF9800 });
    const receiver = new THREE.Mesh(geometry, material);
    receiver.position.set(position.x, position.y, position.z);
    this.scene.add(receiver);
  }

  private createLabels(derivedScene: DerivedScene) {
    const createLabel = (text: string, position: { x: number; y: number; z: number }) => {
      const canvas = document.createElement('canvas');
      const context = canvas.getContext('2d');
      if (!context) return;

      canvas.width = 256;
      canvas.height = 64;
      context.font = '24px Arial';
      context.fillStyle = 'white';
      context.fillText(text, 10, 40);

      const texture = new THREE.CanvasTexture(canvas);
      const spriteMaterial = new THREE.SpriteMaterial({ map: texture });
      const sprite = new THREE.Sprite(spriteMaterial);
      sprite.position.set(position.x, position.y, position.z + 0.5);
      sprite.scale.set(2, 0.5, 1);
      this.scene.add(sprite);
      this.labels.push(sprite);
    };

    createLabel('L', derivedScene.leftSpeaker);
    createLabel('R', derivedScene.rightSpeaker);
    createLabel('RX', derivedScene.receiver);
  }

  private fitCameraToRoom(derivedScene: DerivedScene) {
    const { roomBounds } = derivedScene;
    const maxDim = Math.max(roomBounds.maxX, roomBounds.maxY, roomBounds.maxZ);
    this.camera.position.set(maxDim * 2, maxDim * 2, maxDim * 1.5);
    const target = new THREE.Vector3(roomBounds.maxX / 2, roomBounds.maxY / 2, roomBounds.maxZ / 2);
    this.camera.lookAt(target);
    if (this.controls) {
      this.controls.target.copy(target);
      this.controls.update();
    }
  }

  animate() {
    requestAnimationFrame(() => this.animate());
    if (this.controls) {
      this.controls.update();
    }
    this.renderer.render(this.scene, this.camera);
  }

  dispose() {
    window.removeEventListener('resize', () => this.onResize());
    if (this.controls) {
      this.controls.dispose();
    }
    this.renderer.dispose();
  }

  resetCamera() {
    if (!this.derivedScene || !this.controls) return;
    this.fitCameraToRoom(this.derivedScene);
  }
}