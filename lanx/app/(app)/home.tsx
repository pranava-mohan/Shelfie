import { BookList } from "@/components/BookList";
import { StyleSheet, View, TextInput, TouchableOpacity, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useState } from "react";
import { normalize } from "@/core/normalize";

export default function ExploreScreen() {
  const [search, setSearch] = useState("");
  const [showFilters, setShowFilters] = useState(false);
  const [filters, setFilters] = useState({
    genre: "",
    author: "",
    publisher: "",
  });

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Explore Books</Text>
      </View>
      <View style={styles.searchContainer}>
        <TextInput
          style={styles.input}
          placeholder="Search books..."
          value={search}
          onChangeText={setSearch}
        />
        <TouchableOpacity
          style={styles.filterButton}
          onPress={() => setShowFilters(!showFilters)}
        >
          <Text style={styles.filterButtonText}>
            {showFilters ? "Hide Filters" : "Show Filters"}
          </Text>
        </TouchableOpacity>

        {showFilters && (
          <View style={styles.filtersContainer}>
            <TextInput
              style={styles.input}
              placeholder="Filter by Genre"
              value={filters.genre}
              onChangeText={(text) => setFilters({ ...filters, genre: text })}
            />
            <TextInput
              style={styles.input}
              placeholder="Filter by Author"
              value={filters.author}
              onChangeText={(text) => setFilters({ ...filters, author: text })}
            />
            <TextInput
              style={styles.input}
              placeholder="Filter by Publisher"
              value={filters.publisher}
              onChangeText={(text) => setFilters({ ...filters, publisher: text })}
            />
          </View>
        )}
      </View>
      <BookList filters={{ ...filters, search }} />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#F5F7FA",
  },
  searchContainer: {
    padding: normalize(16),
    backgroundColor: "#F5F7FA",
    borderBottomWidth: 1,
    borderBottomColor: "#e2e8f0",
  },
  header: {
    paddingHorizontal: normalize(16),
    paddingTop: normalize(12),
    paddingBottom: normalize(4),
    backgroundColor: "#F5F7FA",
    alignItems: "center",
  },
  headerTitle: {
    fontSize: normalize(24),
    fontWeight: "bold",
    color: "#1A202C",
  },
  input: {
    backgroundColor: "#f7fafc",
    borderRadius: normalize(8),
    padding: normalize(10),
    marginBottom: normalize(10),
    borderWidth: 1,
    borderColor: "#e2e8f0",
  },
  filterButton: {
    alignSelf: "flex-end",
    padding: normalize(8),
  },
  filterButtonText: {
    color: "#3182ce",
    fontWeight: "600",
  },
  filtersContainer: {
    marginTop: normalize(10),
  },
});

