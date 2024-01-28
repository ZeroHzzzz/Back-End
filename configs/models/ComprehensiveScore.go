package models

type ComprehensiveScore struct {
	userId       int64   `json:"-"`
	academicYear string  `json:"-"`
	totalScore   float64 `json:"-"`
	// 其他综测分数信息字段

	moralQualities        MoralQualities        // 德育素质
	intellectualQualities IntellectualQualities // 智育素质
	physicalQualities     PhysicalQualities     // 体育素质
	aestheticQualities    AestheticQualities    // 美育素质
	laborQualities        LaborQualities        // 劳育素质
	innovationQualities   InnovationQualities   // 创新与实践素质
}

// 德
type MoralQualities struct {
	moralScore             float64                `json:"-"`
	basicAssessmentD1      float64                // 基本评定分D1
	practicalAddSubtractD2 PracticalAddSubtractD2 // 记实加减分D2
}

type PracticalAddSubtractD2 struct {
	collectiveRatingScore       float64            // 集体评定等级分
	socialResponsibilityScore   float64            // 社会责任记实分
	ideologicalLearningScore    float64            // 思政学习加减分
	violationDeduction          float64            // 违纪违规扣分
	honorTitleAddSubtract       float64            // 学生荣誉称号加减分
	honorTitleAddSubtractSource map[string]float64 // 荣誉称号来源
}

// 智
type IntellectualQualities struct {
	intellectualScore float64 `json:"-"`
	averageGPA        float64 // 智育平均学分绩点
}

// 体
type PhysicalQualities struct {
	physicalScore                  float64               `json:"-"`
	PhysicalCourseScoreT1          PhysicalCourseScoreT1 // 体育课程成绩T1
	extracurricularActivityScoreT2 float64               // 课外体育活动成绩T2
}

type PhysicalCourseScoreT1 struct {
	semester1Score float64 // 第一学期得分
	semester2Score float64 // 第二学期得分
}

// 美
type AestheticQualities struct {
	aestheticSorce                      float64            `json:"-"`
	culturalArtPracticeScoreM1          float64            // 文化艺术实践成绩M1
	culturalArtPracticeScoreM1Source    map[string]float64 `json:"-"`
	culturalArtCompetitionScoreM2       float64            // 文化艺术竞赛获奖得分M2
	culturalArtCompetitionScoreM2Source map[string]float64 //竞赛
}

// 劳
type LaborQualities struct {
	laborScore                float64           `json:"-"`
	dailyLaborScoreL1         DailyLaborScoreL1 // 日常劳动分L1
	volunteerServiceScoreL2   float64           // 志愿服务分L2
	internshipTrainingScoreL3 float64           // 实习实训L3
}

type DailyLaborScoreL1 struct {
	basicAssessmentScore    float64 // 寝室日常考核基本分
	civilizedDormitoryScore float64 // "文明寝室"创建、寝室风采展等活动加分
	dormitoryBehaviorScore  float64 // 寝室行为表现与卫生状况加减分
}

// 创新
type InnovationQualities struct {
	innovativeScore                   float64                           `json:"-"`
	innovationEntrepreneurshipScoreC1 InnovationEntrepreneurshipScoreC1 // 创新创业成绩C1
	socialPracticeActivityScoreC2     float64                           // 社会实践活动C2
	socialWorkScoreC3                 float64                           // 社会工作C3
}

type InnovationEntrepreneurshipScoreC1 struct {
	innovationCompetitionScore       float64            // 创新创业竞赛获奖得分
	innovationCompetitionScoreSource map[string]float64 // 竞赛来源
	levelExamScore                   float64            // 水平等级考试
	levelExamScoreSource             map[string]float64 // 考试类型
}
