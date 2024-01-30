package scoredatabase

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 维护数据库关系
func maintainReletion(collection *mongo.Collection) error {
	pipeline := []bson.M{
		{
			"$project": bson.M{
				// 选择包含 "D" 的字段

				"Father": bson.M{
					"$regexMatch": bson.M{
						"input": "$F",
						"regex": "F",
					},
				},

				"dFields": bson.M{
					"$filter": bson.M{
						"input": bson.M{"$objectToArray": "$$ROOT"},
						"as":    "field",
						"cond": bson.M{
							"$regexMatch": bson.M{
								"input": "$$field.k",
								"regex": "d",
							},
						},
					},
				},
				// 获取嵌套文档 "L" 的第一个字段的值
				"LFirstField": bson.M{"$first": "$L"},
			},
		},
		{
			"$set": bson.M{
				"Father": bson.M{
					"$sum": bson.A{
						// 求和 "dFields" 中的所有值
						bson.M{
							"$reduce": bson.M{
								"input":        "$dFields",
								"initialValue": 0,
								"in":           bson.M{"$sum": bson.A{"$$value", "$$this.v"}},
							},
						},
						// 求和 "LFirstField"
						bson.M{
							"$ifNull": bson.A{"$LFirstField", 0},
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				// 选择包含 "D" 的字段
				"Total": 1,
				"Father": bson.M{
					"$regexMatch": bson.M{
						"input": "$F",
						"regex": "F",
					},
				},

				"DFields": bson.M{
					"$filter": bson.M{
						"input": bson.M{"$objectToArray": "$$ROOT"},
						"as":    "field",
						"cond": bson.M{
							"$regexMatch": bson.M{
								"input": "$$field.k",
								"regex": "D",
							},
						},
					},
				},
			},
		},
		{
			"$set": bson.M{
				"Total": bson.M{
					"$sum": bson.A{
						"$Father",
						bson.M{
							"$reduce": bson.M{
								"input":        "$DFields",
								"initialValue": 0,
								"in":           "$$value + $$this.v",
							},
						},
					},
				},
			},
		},
	}
	// 聚合
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Println("Error aggregating documents:", err)
		return err
	}
	defer cursor.Close(context.TODO())
	// 处理聚合结果
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			fmt.Println("Error decoding document:", err)
			return err
		}

		// // 处理 result，可以输出到控制台或进行其他操作
		// fmt.Println(result)
		// return result, nil
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Error iterating cursor:", err)
		return err
	}
	return nil
}
