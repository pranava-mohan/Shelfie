import React, { useState, useEffect } from "react";
import { Text, View, StyleSheet, Button, Linking, TextInput } from "react-native";
import { CameraView, Camera, BarcodeScanningResult } from "expo-camera";
import { useClient } from "@/core/hooks/useClient";
import Toast from "react-native-toast-message";
import { useAppSelector } from "@/core/store/store";
import { useRouter } from "expo-router";
import { API } from "@/config";

export default function ChekInPage() {
  const [hasPermission, setHasPermission] = useState(false);
  const [scanned, setScanned] = useState(false);
  const { client } = useClient();
  const router = useRouter();

  const userid = useAppSelector((state) => state.auth.id);

  useEffect(() => {
    const getCameraPermissions = async () => {
      const { status } = await Camera.requestCameraPermissionsAsync();
      setHasPermission(status === "granted");
    };

    getCameraPermissions();
  }, []);

  const handleBarCodeScanned = async ({
    type,
    data,
  }: { type: string, data: string, cornerPoints?: any[], bounds?: any }) => {
    if (data.length == 256) {
      setScanned(true);
      await client.post(`${API.BASE_URL}/check-in/` + data, {
        UserID: userid,
      });
      Toast.show({
        type: "success",
        text1: "Auth Successful",
      });
      setScanned(false);
      router.replace("/(app)/home");
    }
  };

  if (hasPermission === null) {
    return <Text>Requesting for camera permission</Text>;
  }
  if (hasPermission === false) {
    return <Text>No access to camera</Text>;
  }

  return (
    <View style={styles.container}>
      <CameraView
        onBarcodeScanned={scanned ? undefined : handleBarCodeScanned}
        barcodeScannerSettings={{
          barcodeTypes: ["qr"],
        }}
        style={styles.camera}
      />
      <View style={styles.overlay}>
        <View style={styles.unfocusedContainer}></View>
        <View style={styles.middleContainer}>
          <View style={styles.unfocusedContainer}></View>
          <View style={styles.focusedContainer}></View>
          <View style={styles.unfocusedContainer}></View>
        </View>
        <View style={styles.unfocusedContainer}></View>
      </View>
      {__DEV__ && (
        <View style={{ position: "absolute", bottom: 50, left: 20, right: 20, backgroundColor: 'white', padding: 10, borderRadius: 10 }}>
          <TextInput
            placeholder="Enter Room ID"
            style={{ borderWidth: 1, padding: 10, marginBottom: 10 }}
            onSubmitEditing={(e) => handleBarCodeScanned({ type: "qr", data: e.nativeEvent.text })}
          />
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: "column",
    justifyContent: "center",
  },
  camera: {
    ...StyleSheet.absoluteFillObject,
  },
  overlay: {
    position: "absolute",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  unfocusedContainer: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.5)",
  },
  middleContainer: {
    flexDirection: "row",
    flex: 1,
  },
  focusedContainer: {
    flex: 4,
  },
});
