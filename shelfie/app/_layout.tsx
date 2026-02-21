import { DefaultTheme, ThemeProvider } from "@react-navigation/native";
import { Slot, useRouter, SplashScreen } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useEffect } from "react";
import { Provider } from "react-redux";
import { store, useAppSelector } from "@/core/store/store";
import { useAuth } from "@/core/hooks/useAuth";
import Toast from "react-native-toast-message";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

export default function RootLayout() {
  return (
    <Provider store={store}>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider value={DefaultTheme}>
          <StatusBar style="auto" />
          <AuthStack />
        </ThemeProvider>
      </QueryClientProvider>
      <Toast />
    </Provider>
  );
}

const AuthStack = () => {
  const { isLoading } = useAppSelector((state) => state.auth);
  const router = useRouter();
  const { check, login, logout } = useAuth();

  useEffect(() => {
    const checkToken = async () => {
      const [status, token] = await check();
      if (status) {
        await login(token as string);
        router.replace("/(app)/home");
      } else {
        await logout();
        router.replace("/(auth)");
      }
    };

    checkToken();
  }, []);

  useEffect(() => {
    if (!isLoading) {
      SplashScreen.hideAsync();
    }
  }, [isLoading]);
  if (isLoading) {
    return null;
  }

  return <Slot />;
};
