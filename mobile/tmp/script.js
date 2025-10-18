// DOM要素の取得
const selector = document.getElementById('template-selector');
const svgContainer = document.getElementById('preview-svg');
const guideCanvas = document.getElementById('guide-canvas');
const guideCtx = guideCanvas.getContext('2d');
const cameraView = document.getElementById('camera-view');
const startButton = document.getElementById('start-camera');
const takePhotoButton = document.getElementById('take-photo');
const restartButton = document.getElementById('restart');
const resultContainer = document.getElementById('collage-result-container');
const colors = ['rgba(255, 0, 0, 0.3)', 'rgba(0, 0, 255, 0.3)', 'rgba(0, 255, 0, 0.3)', 'rgba(255, 255, 0, 0.3)', 'rgba(0, 255, 255, 0.3)'];
let templates = [];
let capturedImages = [];
let currentFrameIndex = 0;
let cameraStream = null;

// JSONファイルを読み込んで処理を実行
fetch('templates.json')
    .then(response => response.json())
    .then(data => {
        templates = data;
        if (templates.length === 0) return;

        // プルダウンメニューを作成
        templates.forEach((template, index) => {
            const option = document.createElement('option');
            option.value = index;
            option.textContent = template.name;
            selector.appendChild(option);
        });

        // 最初のテンプレートを描画
        renderTemplate(templates[0]);
        renderGuide(templates[0]);

        // プルダウンが変更されたら、対応するテンプレートを描画
        selector.addEventListener('change', (event) => {
            const selectedIndex = event.target.value;
            renderTemplate(templates[selectedIndex]);
            resetShootingProgress();
        });
    })
    .catch(error => console.error('Error loading templates:', error));

// SVGを描画する関数
function renderTemplate(template) {
    svgContainer.innerHTML = '';
    if (template.viewBox) {
        svgContainer.setAttribute('viewBox', template.viewBox);
    }
    template.frames.forEach((frame, index) => {
        const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        path.setAttribute('d', frame.path);
        path.setAttribute('fill', colors[index % colors.length]);
        svgContainer.appendChild(path);
    });
}

// ガイドを描画する関数
function renderGuide(template, activeFrameIndex = -1) {
    guideCanvas.width = 500;
    guideCanvas.height = 500;
    guideCtx.clearRect(0, 0, guideCanvas.width, guideCanvas.height);

    if (activeFrameIndex !== -1 && cameraStream) {
        // 背景を薄暗くする
        guideCtx.fillStyle = 'rgba(0, 0, 0, 0.5)';
        guideCtx.fillRect(0, 0, guideCanvas.width, guideCanvas.height);

        // アクティブなフレーム部分だけをクリアする
        const activeFrame = template.frames[activeFrameIndex];
        const path = new Path2D(activeFrame.path);
        const matrix = new DOMMatrix().scale(guideCanvas.width, guideCanvas.height);
        const transformedPath = new Path2D();
        transformedPath.addPath(path, matrix);
        guideCtx.save();
        guideCtx.clip(transformedPath);
        guideCtx.clearRect(0, 0, guideCanvas.width, guideCanvas.height);
        guideCtx.restore();
    }

    // すべてのフレームの枠線を描画
    guideCtx.strokeStyle = 'white';
    guideCtx.lineWidth = 2;
    template.frames.forEach(frame => {
        const path = new Path2D(frame.path);
        const matrix = new DOMMatrix().scale(guideCanvas.width, guideCanvas.height);
        const transformedPath = new Path2D();
        transformedPath.addPath(path, matrix);
        guideCtx.stroke(transformedPath);
    });
}

// SVGパスを解析してバウンディングボックスを取得する関数
function getFrameBoundingBox(frame) {
    const path = frame.path;
    const points = [];
    let currentX = 0;
    let currentY = 0;

    const regex = /([MLHVZ])([^MLHVZ]*)/gi;
    let match;

    while ((match = regex.exec(path)) !== null) {
        const command = match[1].toUpperCase();
        const args = match[2].trim().split(/[\s,]+/).map(parseFloat).filter(n => !isNaN(n));

        if (command === 'M' || command === 'L') {
            for (let i = 0; i < args.length; i += 2) {
                currentX = args[i];
                currentY = args[i + 1];
                points.push([currentX, currentY]);
            }
        } else if (command === 'H') {
            for (let i = 0; i < args.length; i++) {
                currentX = args[i];
                points.push([currentX, currentY]);
            }
        } else if (command === 'V') {
            for (let i = 0; i < args.length; i++) {
                currentY = args[i];
                points.push([currentX, currentY]);
            }
        } else if (command === 'Z') {
            // Path closes
        }
    }

    if (points.length === 0) {
        return { x: 0, y: 0, width: 0, height: 0 };
    }

    const minX = Math.min(...points.map(p => p[0]));
    const minY = Math.min(...points.map(p => p[1]));
    const maxX = Math.max(...points.map(p => p[0]));
    const maxY = Math.max(...points.map(p => p[1]));

    return { x: minX, y: minY, width: maxX - minX, height: maxY - minY };
}

// video要素内の映像の表示スケールとオフセットを計算する関数
function getVideoDisplayInfo(videoElement) {
    if (!videoElement || !videoElement.videoWidth || videoElement.videoWidth === 0) {
        return { scale: 1, offsetX: 0, offsetY: 0 };
    }
    const videoRatio = videoElement.videoWidth / videoElement.videoHeight;
    const elementRatio = videoElement.clientWidth / videoElement.clientHeight;
    let scale = 1;
    let offsetX = 0;
    let offsetY = 0;

    if (videoRatio > elementRatio) { // 映像が要素より横長 -> 上下に黒帯
        scale = videoElement.clientWidth / videoElement.videoWidth;
        offsetY = (videoElement.clientHeight - (videoElement.videoHeight * scale)) / 2;
    } else { // 映像が要素より縦長 -> 左右に黒帯
        scale = videoElement.clientHeight / videoElement.videoHeight;
        offsetX = (videoElement.clientWidth - (videoElement.videoWidth * scale)) / 2;
    }
    return { scale, offsetX, offsetY };
}


// コラージュを作成する関数
function createCollage(images) {
    const template = templates[selector.value];
    const collageCanvas = document.createElement('canvas');
    const canvasSize = 500;
    collageCanvas.width = canvasSize;
    collageCanvas.height = canvasSize;
    const collageCtx = collageCanvas.getContext('2d');

    collageCtx.fillStyle = 'white';
    collageCtx.fillRect(0, 0, collageCanvas.width, collageCanvas.height);

    const videoInfo = getVideoDisplayInfo(cameraView);

    images.forEach((image, index) => {
        const frame = template.frames[index];
        const frameBbox = getFrameBoundingBox(frame);

        // 1. ガイドの領域をピクセル単位に変換 (on 500x500 view)
        const viewBoxPx = {
            x: frameBbox.x * canvasSize,
            y: frameBbox.y * canvasSize,
            width: frameBbox.width * canvasSize,
            height: frameBbox.height * canvasSize
        };

        // 2. ガイド領域をソース画像上の座標に変換 (レターボックスを考慮)
        const sx = (viewBoxPx.x - videoInfo.offsetX) / videoInfo.scale;
        const sy = (viewBoxPx.y - videoInfo.offsetY) / videoInfo.scale;
        const sWidth = viewBoxPx.width / videoInfo.scale;
        const sHeight = viewBoxPx.height / videoInfo.scale;

        // 3. ソース画像からはみ出ないようにクリッピング
        const clip_sx = Math.max(0, sx);
        const clip_sy = Math.max(0, sy);
        const right_sx = Math.min(image.width, sx + sWidth);
        const bottom_sy = Math.min(image.height, sy + sHeight);
        const clip_sWidth = right_sx - clip_sx;
        const clip_sHeight = bottom_sy - clip_sy;

        if (clip_sWidth <= 0 || clip_sHeight <= 0) {
            return; // ソース領域がない場合はスキップ
        }
        
        // 4. コラージュキャンバス上の描画先領域を定義
        const path = new Path2D(frame.path);
        const matrix = new DOMMatrix().scale(canvasSize, canvasSize);
        const transformedPath = new Path2D();
        transformedPath.addPath(path, matrix);

        collageCtx.save();
        collageCtx.clip(transformedPath);

        // 5. ソースの比率を維持して描画先を計算 (cover)
        const dx = viewBoxPx.x;
        const dy = viewBoxPx.y;
        const dWidth = viewBoxPx.width;
        const dHeight = viewBoxPx.height;

        const sourceAspectRatio = clip_sWidth / clip_sHeight;
        const destAspectRatio = dWidth / dHeight;

        let final_dx = dx;
        let final_dy = dy;
        let final_dWidth = dWidth;
        let final_dHeight = dHeight;

        if (sourceAspectRatio > destAspectRatio) {
            final_dHeight = dHeight;
            final_dWidth = final_dHeight * sourceAspectRatio;
            final_dx = dx - (final_dWidth - dWidth) / 2;
        } else {
            final_dWidth = dWidth;
            final_dHeight = final_dWidth / sourceAspectRatio;
            final_dy = dy - (final_dHeight - dHeight) / 2;
        }

        collageCtx.drawImage(image,
            clip_sx, clip_sy, clip_sWidth, clip_sHeight,
            final_dx, final_dy, final_dWidth, final_dHeight
        );

        collageCtx.restore();
    });

    resultContainer.innerHTML = '';
    resultContainer.appendChild(collageCanvas);
    takePhotoButton.disabled = true;
}


// カメラを停止する関数
function stopCamera() {
    if (cameraStream) {
        cameraStream.getTracks().forEach(track => track.stop());
        cameraStream = null;
        cameraView.srcObject = null;
    }
}

// 撮影プログレスをリセットする関数
function resetShootingProgress() {
    capturedImages = [];
    currentFrameIndex = 0;
    resultContainer.innerHTML = '';
    takePhotoButton.disabled = !cameraStream;

    if (cameraStream) {
        const template = templates[selector.value];
        renderGuide(template, currentFrameIndex);
    }
}

// リスタート処理（完全リセット）
function restart() {
    stopCamera();
    capturedImages = [];
    currentFrameIndex = 0;
    resultContainer.innerHTML = '';
    takePhotoButton.disabled = true;
    restartButton.disabled = true;
    startButton.disabled = false;
    guideCtx.clearRect(0, 0, guideCanvas.width, guideCanvas.height);
    const template = templates[selector.value];
    if (template) {
        renderGuide(template);
    }
}

// --- イベントリスナー ---
startButton.addEventListener('click', async () => {
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true });
        cameraStream = stream;
        cameraView.srcObject = stream;
        // video要素のメタデータが読み込まれるのを待つ
        cameraView.onloadedmetadata = () => {
            takePhotoButton.disabled = false;
            restartButton.disabled = false;
            startButton.disabled = true;
            const template = templates[selector.value];
            renderGuide(template, currentFrameIndex);
        };
    } catch (error) {
        console.error('Error accessing camera:', error);
        alert('カメラにアクセスできませんでした。他のアプリでカメラを使用していないか、またはカメラが正しく接続されているか確認してください。');
    }
});

takePhotoButton.addEventListener('click', () => {
    const template = templates[selector.value];
    if (currentFrameIndex < template.frames.length) {
        const tempCanvas = document.createElement('canvas');
        tempCanvas.width = cameraView.videoWidth;
        tempCanvas.height = cameraView.videoHeight;
        const tempCtx = tempCanvas.getContext('2d');
        tempCtx.drawImage(cameraView, 0, 0);
        
        const image = new Image();
        const imageLoadPromise = new Promise(resolve => {
            image.onload = () => resolve(image);
        });
        image.src = tempCanvas.toDataURL('image/png');
        capturedImages.push(imageLoadPromise);

        currentFrameIndex++;

        if (currentFrameIndex < template.frames.length) {
            renderGuide(template, currentFrameIndex);
        } else {
            Promise.all(capturedImages).then(loadedImages => {
                createCollage(loadedImages);
                guideCtx.clearRect(0, 0, guideCanvas.width, guideCanvas.height);
            });
        }
    }
});

restartButton.addEventListener('click', restart);
