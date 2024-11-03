import { integer, pgTable, timestamp, uuid, varchar } from "drizzle-orm/pg-core";

export const usersTable = pgTable("users", {
  id: uuid().primaryKey().defaultRandom(),
  name: varchar({ length: 255 }).notNull(),
  email: varchar({ length: 255 }).notNull().unique(),
  age: integer().notNull(),
  imageUrl: varchar({ length: 255 }),
  stripeCustomerId: varchar({ length: 255 }).notNull().unique(),
  createdAt: timestamp().defaultNow(),
});

export const plansTable = pgTable("plans", {
  id: uuid().primaryKey().defaultRandom(),
  name: varchar().notNull().unique(),
  description: varchar().notNull(),
  // priceCents: integer().notNull(), // TODO: this would be different depending on country
  stripeProductId: varchar().notNull().unique(),
  createdAt: timestamp().defaultNow(),
});
