package grpc

// func TestError_WithRoot(t *testing.T) {
// 	type fields struct {
// 		Root   xerror.Root
// 		Status status.Status
// 	}
// 	type args struct {
// 		r *xerror.Root
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *Error
// 	}{
// 		{},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			st := status.New(codes.Unknown, "something unknown happened")
// 			infoDetail := errdetails.ErrorInfo{
// 				Reason: "some reason",
// 				Domain: "some domain",
// 				Metadata: map[string]string{
// 					"key1": "value1",
// 				},
// 			}

// 			// anypb, err := anypb.New(&infoDetail)
// 			// require.NoError(t, err, "failed to create anypb")

// 			// st, err = st.WithDetails(anypb)
// 			// require.NoError(t, err, "failed to attach details to status")

// 			var err error
// 			st, err = st.WithDetails(&infoDetail)
// 			require.NoError(t, err, "failed to attach details to status")

// 			b, err := protojson.Marshal(st.Proto())
// 			require.NoError(t, err, "failed to marshal status to json")

// 			t.Logf("error as json:\n%s", string(b))
// 		})
// 	}
// }
