package utils

// func TestCalculateAQI(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		concentration float64
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want int
// 	}{
// 		{
// 			name: "Good",
// 			args: args{
// 				concentration: 0.0,
// 			},
// 			want: 0,
// 		},
// 		{
// 			name: "Moderate",
// 			args: args{
// 				concentration: 12.1,
// 			},
// 			want: 51,
// 		},
// 		{
// 			name: "Unhealthy for Sensitive Groups",
// 			args: args{
// 				concentration: 35.5,
// 			},
// 			want: 101,
// 		},
// 		{
// 			name: "Unhealthy",
// 			args: args{
// 				concentration: 55.5,
// 			},
// 			want: 151,
// 		},
// 		{
// 			name: "Very Unhealthy",
// 			args: args{
// 				concentration: 150.5,
// 			},
// 			want: 201,
// 		},
// 		{
// 			name: "Hazardous",
// 			args: args{
// 				concentration: 250.5,
// 			},
// 			want: 301,
// 		},
// 		{
// 			name: "Hazardous",
// 			args: args{
// 				concentration: 350.5,
// 			},
// 			want: 401,
// 		},
// 		{
// 			name: "Hazardous",
// 			args: args{
// 				concentration: 500.5,
// 			},
// 			want: 501,
// 		},
// 		{
// 			name: "Hazardous",
// 			args: args{
// 				concentration: 1000.5,
// 			},
// 			want: 501,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := CalculateAQI(tt.args.concentration); got != tt.want {
// 				t.Errorf("CalculateAQI() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
