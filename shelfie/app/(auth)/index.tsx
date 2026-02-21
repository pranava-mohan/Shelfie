// index.tsx
import { normalize } from "@/core/normalize";
import { Ionicons } from "@expo/vector-icons";
import { useRouter } from "expo-router";
import { Pressable, Text, View } from "react-native";
import * as Linking from "expo-linking";
import * as WebBrowser from "expo-web-browser";
import { useEffect } from "react";
import { useAuth } from "@/core/hooks/useAuth";
import { AUTH } from "@/config";

export default function Index() {
  const router = useRouter();
  const { login } = useAuth();

  const handleDeepLink = async (event: any) => {
    let { queryParams } = Linking.parse(event.url);

    if (queryParams?.token) {
      await login(queryParams?.token as string);
      router.replace("/(app)/home");
    }
  };

  useEffect(() => {
    const subscription = Linking.addEventListener("url", handleDeepLink);

    Linking.getInitialURL().then((url) => {
      if (url) {
        handleDeepLink({ url });
      }
    });

    return () => {
      subscription.remove();
    };
  }, []);

  const handleDauthLogin = async () => {
    try {
      const redirectUrl = Linking.createURL("auth");
      if (!AUTH.GOOGLE_URL) throw new Error("Google Auth URL is not defined");
      await WebBrowser.openAuthSessionAsync(AUTH.GOOGLE_URL, redirectUrl);
    } catch (e: any) {
      alert("Error opening browser: " + (e.message || String(e)));
    }
  };

  return (
    <View
      style={{
        flex: 1,
        alignItems: "center",
        justifyContent: "space-between",
        paddingVertical: 40,
      }}
    >
      <View
        style={{
          flex: 1,
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Ionicons name="book-sharp" size={normalize(100)} color="black" />
      </View>

      <Pressable
        style={{
          backgroundColor: "#000000",
          paddingVertical: 18,
          paddingHorizontal: 80,
          borderRadius: 10,
          marginBottom: 20,
        }}
        onPress={() => {
          handleDauthLogin();
        }}
      >
        <Text
          style={{
            fontSize: normalize(20),
            color: "#FFFFFF",
          }}
        >
          Get Started
        </Text>
      </Pressable>
    </View>
  );
}
