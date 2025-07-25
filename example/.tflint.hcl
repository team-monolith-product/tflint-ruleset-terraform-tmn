plugin "terraform-sort" {
    enabled = true
}

plugin "terraform" {
    enabled = false
}

rule "terraform_list_order" {
    enabled = true
}


rule "terraform_variables_order" {
    enabled = true
    group_required = true
}