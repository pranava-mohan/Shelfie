import { Dimensions, PixelRatio, Platform } from "react-native";

const { width: SCREEN_WIDTH } = Dimensions.get("window");

const scale = SCREEN_WIDTH / 375;

export function normalize(size: number) {
  const newSize = size * scale;
  if (Platform.OS === "ios") {
    return Math.round(PixelRatio.roundToNearestPixel(newSize));
  } else {
    return Math.round(PixelRatio.roundToNearestPixel(newSize)) - 2;
  }
}
