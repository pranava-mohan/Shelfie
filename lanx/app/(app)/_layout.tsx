import { Tabs } from "expo-router";
import React from "react";

import { Ionicons } from "@expo/vector-icons";
import { normalize } from "@/core/normalize";

export default function TabLayout() {
  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: "#000000",
        headerShown: false,
      }}
    >
      <Tabs.Screen
        name="scan"
        options={{
          title: "Scan",
          tabBarIcon: ({ color }) => (
            <Ionicons size={28} name="qr-code" color={color} />
          ),
          tabBarLabelStyle: {
            fontSize: normalize(14),
          },
        }}
      />
      <Tabs.Screen
        name="home"
        options={{
          title: "Catalog",
          tabBarIcon: ({ color }) => (
            <Ionicons size={28} name="list" color={color} />
          ),
          tabBarLabelStyle: {
            fontSize: normalize(14),
          },
        }}
      />
      <Tabs.Screen
        name="my_books"
        options={{
          title: "My Books",
          tabBarIcon: ({ color }) => (
            <Ionicons size={28} name="book" color={color} />
          ),
          tabBarLabelStyle: {
            fontSize: normalize(14),
          },
        }}
      />
    </Tabs>
  );
}
