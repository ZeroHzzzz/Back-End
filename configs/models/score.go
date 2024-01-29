package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"` //文档id
	UserId                string             `bson:"user_id"`
	MoralQualities        int                `bson:"德育素质"`
	D1                    int                `bson:"基本评定分D1"`
	D2                    int                `bson:"记实加减分D2"`
	D2_1                  int                `bson:"集体评定等级分"`
	D2_2                  int                `bson:"社会责任记实分"`
	D2_3                  int                `bson:"思政学习加减分"`
	D2_4                  int                `bson:"违纪违规扣分"`
	D2_5                  Score_Source       `bson:"学生荣誉称号加减分"`
	IntellectualQualities int                `bson:"智育素质"` //因为这一项下面就一个分项，就不单独罗列
	PhysicalQualities     int                `bson:"体育素质"`
	T1                    int                `bson:"体育课程成绩"`
	T2                    int                `bson:"课外体育活动成绩"`
	T2_1                  Score_Source       `bson:"体育竞赛获奖得分"`
	T2_2                  int                `bson:"早锻炼得分"`
	AestheticQualities    int                `bson:"美育素质"`
	M1                    int                `bson:"文化艺术实践成绩"`
	M2                    Score_Source       `bson:"文化艺术竞赛获奖得分"`
	LaborQualities        int                `bson:"劳育素质"`
	L1                    int                `bson:"寝室日常考核基本分"`
	L1_1                  int                `bson:"寝室日常考核基本分"`
	L1_2                  int                `bson:"“文明寝室”创建、寝室风采展等活动加分"`
	L1_3                  int                `bson:"寝室行为表现与卫生状况加减分"`
	L2                    int                `bson:"志愿服务分"`
	L3                    int                `bson:"实习实训"`
	InnovationQualities   int                `bson:"创新与实践素质"`
	C1                    int                `bson:"创新创业成绩"`
	C1_1                  Score_Source       `bson:"创新创业竞赛获奖得分"`
	C1_2                  Score_Source       `bson:"水平等级考试"`
	C2                    Score_Source       `bson:"社会实践活动"`
	C3                    int                `bson:"社会工作"`
}

// 这个类型是用来记录来源的
type Score_Source struct {
	Score  int      `bson:"sorce"`
	Source []string `bson:"Source"`
}
