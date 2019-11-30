package labx_subs

type Resolver struct{}

func (r *Resolver) Access() AccessResolver {
	return &accessResolver{r}
}
func (r *Resolver) Supplier() SupplierResolver {
	return &supplierResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
