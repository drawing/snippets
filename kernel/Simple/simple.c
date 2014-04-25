#include <linux/kernel.h>
#include <linux/types.h>
#include <linux/init.h>
#include <linux/module.h>
#include <linux/slab.h>
#include <linux/rbtree.h>

struct simple
{
	int data;
	struct rb_node node;
};

struct simple * search_simple(struct rb_root * root, int data)
{
	struct simple * p;
	struct rb_node * pnode = root->rb_node;
	while (pnode)
	{
		p = rb_entry(pnode, struct simple, node);
		if (data < p->data) {
			pnode = pnode->rb_left;
		} else if (data > p->data){
			pnode = pnode->rb_right;
		} else {
			return p;
		}
	}
	return NULL;
}

struct simple * insert_simple(struct rb_root * root,
		int data, struct simple * new)
{
	struct simple * p = NULL;
	struct rb_node ** pnode = &root->rb_node;
	struct rb_node * parent = NULL;
	while (*pnode)
	{
		parent = *pnode;
		p = rb_entry(parent, struct simple, node);
		if (data < p->data) {
			pnode = &(*pnode)->rb_left;
		} else if (data > p->data){
			pnode = &(*pnode)->rb_right;
		} else {
			return p;
		}
	}
	rb_link_node(&new->node, parent, pnode);
	rb_insert_color(&new->node, root);
	return p;
}

struct rb_root g_root = RB_ROOT;

static int __init initialization(void)
{
	int i;
	struct simple * p;
	for (i = 0; i < 10; i++) {
		p = (struct simple *)kmalloc(sizeof (struct simple),
				GFP_KERNEL);
		p->data = i * 10;
		rb_init_node(&p->node);
		insert_simple(&g_root, i * 10, p);
		printk(KERN_INFO "insert %d\n", i * 10);
	}
	p = search_simple(&g_root, 30);
	if (p != NULL) {
		printk(KERN_INFO "find it\n");
	}

	return 0;
}

static void __exit cleanup(void)
{
	struct simple * p;
	struct rb_node * pnode = rb_first(&g_root);
	while (pnode != NULL) {
		p = rb_entry(pnode, struct simple, node);		
		printk(KERN_INFO "delete %d\n", p->data);
		rb_erase(pnode, &g_root);
		kfree(p);
		pnode = rb_first(&g_root);
	}
}

module_init(initialization);
module_exit(cleanup);

MODULE_AUTHOR("cppbreak cppbreak@gmail.com");
MODULE_DESCRIPTION("A simple linux kernel module");
MODULE_VERSION("V0.1");
MODULE_LICENSE("Dual BSD/GPL");

