import cv2
import pytesseract
import numpy as np


def get_grayscale(image):
    return cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)


def remove_noise(image):
    return cv2.medianBlur(image, 5)


def deskew(image):
    coords = np.column_stack(np.where(image > 0))
    angle = cv2.minAreaRect(coords)[-1]
    if angle < -45:
        angle = -(90 + angle)
    else:
        angle = -angle
    (h, w) = image.shape[:2]
    center = (w // 2, h // 2)
    M = cv2.getRotationMatrix2D(center, angle, 1.0)
    rotated = cv2.warpAffine(
        image, M, (w, h), flags=cv2.INTER_CUBIC, borderMode=cv2.BORDER_REPLICATE
    )
    return rotated


def main(path):
    img = cv2.imread(path)

    img = get_grayscale(img)
    img = remove_noise(img)
    img = deskew(img)

    # h, w = img.shape
    # boxes = pytesseract.image_to_boxes(img)
    # for b in boxes.splitlines():
    #     b = b.split(" ")
    #     img = cv2.rectangle(
    #         img, (int(b[1]), h - int(b[2])), (int(b[3]), h - int(b[4])), (0, 255, 0), 2
    #     )
    # cv2.imshow("img", img)
    # cv2.waitKey(0)

    text = pytesseract.image_to_string(img, lang="jpn", config="--psm 12")
    print(text)
