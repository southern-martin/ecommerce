import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Badge } from '@/shared/components/ui/badge';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog';
import {
  ArrowLeft,
  Package,
  Tags,
  Settings2,
  Loader2,
  Plus,
  X,
  Trash2,
  Wand2,
  Save,
  Check,
  ImageIcon,
} from 'lucide-react';
import { ImageUploader } from '@/shared/components/forms/ImageUploader';
import {
  useAdminProduct,
  useAdminUpdateProduct,
  useAdminProductOptions,
  useAdminAddOption,
  useAdminRemoveOption,
  useAdminProductVariants,
  useAdminGenerateVariants,
  useAdminUpdateVariant,
  useAdminProductAttributes,
  useAdminSetProductAttributes,
} from '../hooks/useAdminProductMgmt';
import {
  useCategories,
  useAttributeGroups,
  useGroupAttributes,
} from '../hooks/useAdminProducts';
import type { Attribute } from '../services/admin-product.api';

export default function AdminProductEditPage() {
  const { id } = useParams<{ id: string }>();
  const { data: product, isLoading } = useAdminProduct(id || '');
  const updateProduct = useAdminUpdateProduct();
  const { data: categories } = useCategories();
  const { data: attributeGroups } = useAttributeGroups();

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-[500px] w-full" />
      </div>
    );
  }

  if (!product || !id) {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <p className="text-muted-foreground">Product not found</p>
        <Button asChild variant="link" className="mt-2">
          <Link to="/admin/products">Back to Products</Link>
        </Button>
      </div>
    );
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const p = product as any;
  const productType = p._product_type || p.product_type || 'simple';
  const isConfigurable = productType === 'configurable';
  const attributeGroupId = p._attribute_group_id || '';

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button asChild variant="ghost" size="icon">
          <Link to="/admin/products">
            <ArrowLeft className="h-5 w-5" />
          </Link>
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">{product.name}</h1>
          <p className="text-sm text-muted-foreground">
            {product.slug} &middot;{' '}
            <Badge variant={isConfigurable ? 'default' : 'outline'} className="text-xs">
              {isConfigurable ? 'Configurable' : 'Simple'}
            </Badge>
          </p>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="general">
        <TabsList>
          <TabsTrigger value="general" className="gap-1.5">
            <Package className="h-4 w-4" />
            General
          </TabsTrigger>
          <TabsTrigger value="attributes" className="gap-1.5">
            <Tags className="h-4 w-4" />
            Attributes
          </TabsTrigger>
          {isConfigurable && (
            <TabsTrigger value="variations" className="gap-1.5">
              <Settings2 className="h-4 w-4" />
              Variations
            </TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="general">
          <GeneralTab
            productId={id}
            product={product}
            categories={categories || []}
            attributeGroups={attributeGroups || []}
            onUpdate={updateProduct}
            currentAttributeGroupId={attributeGroupId}
          />
        </TabsContent>

        <TabsContent value="attributes">
          <AttributesTab
            productId={id}
            attributeGroupId={attributeGroupId}
            isConfigurable={isConfigurable}
          />
        </TabsContent>

        {isConfigurable && (
          <TabsContent value="variations">
            <VariationsTab productId={id} product={product} attributeGroupId={attributeGroupId} />
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}

// ── General Tab ──
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function GeneralTab({ productId, product, categories, attributeGroups, onUpdate, currentAttributeGroupId }: any) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const p = product as any;
  const [name, setName] = useState(product.name);
  const [description, setDescription] = useState(product.description);
  const [priceCents, setPriceCents] = useState(product.price);
  const [currency, setCurrency] = useState(p._currency || 'USD');
  const [status, setStatus] = useState(p._status || 'draft');
  const [categoryId, setCategoryId] = useState(product.category?.id || '');
  const [attributeGroupId, setAttributeGroupId] = useState(currentAttributeGroupId);
  const [stockQuantity, setStockQuantity] = useState(product.stock_quantity || 0);
  const [imageUrls, setImageUrls] = useState<string[]>(product.images?.map((img: { url: string }) => img.url) || []);
  const [saveMessage, setSaveMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleSave = () => {
    setSaveMessage(null);
    onUpdate.mutate(
      {
        id: productId,
        data: {
          name,
          description,
          base_price_cents: priceCents,
          currency,
          status,
          category_id: categoryId,
          attribute_group_id: attributeGroupId,
          tags: p._tags || [],
          image_urls: imageUrls,
        },
      },
      {
        onSuccess: () => setSaveMessage({ type: 'success', text: 'Product saved successfully.' }),
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        onError: (err: any) => setSaveMessage({ type: 'error', text: err?.response?.data?.error || err?.message || 'Failed to save product.' }),
      }
    );
  };

  const productType = p._product_type || p.product_type || 'simple';

  return (
    <Card>
      <CardHeader>
        <CardTitle>General Information</CardTitle>
        <CardDescription>Basic product details, pricing, and categorization</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="name">Product Name</Label>
            <Input id="name" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <select
              id="status"
              value={status}
              onChange={(e) => setStatus(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={4}
            className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>

        <div className="grid gap-4 sm:grid-cols-3">
          <div className="space-y-2">
            <Label htmlFor="price">Price (cents)</Label>
            <Input
              id="price"
              type="number"
              value={priceCents}
              onChange={(e) => setPriceCents(Number(e.target.value))}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="currency">Currency</Label>
            <select
              id="currency"
              value={currency}
              onChange={(e) => setCurrency(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
              <option value="GBP">GBP</option>
            </select>
          </div>
          {productType === 'simple' && (
            <div className="space-y-2">
              <Label htmlFor="stock">Stock Quantity</Label>
              <Input
                id="stock"
                type="number"
                value={stockQuantity}
                onChange={(e) => setStockQuantity(Number(e.target.value))}
              />
            </div>
          )}
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="category">Category</Label>
            <select
              id="category"
              value={categoryId}
              onChange={(e) => setCategoryId(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="">Select a category...</option>
              {categories.map((cat: { id: string; name: string }) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="attributeGroup">Attribute Group</Label>
            <select
              id="attributeGroup"
              value={attributeGroupId}
              onChange={(e) => setAttributeGroupId(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="">No attribute group</option>
              {attributeGroups.map((g: { id: string; name: string }) => (
                <option key={g.id} value={g.id}>
                  {g.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Product Images */}
        <div className="space-y-2">
          <Label>Product Images</Label>
          <ImageUploader
            value={imageUrls}
            onChange={setImageUrls}
            maxFiles={8}
            maxSizeMB={5}
            ownerType="product"
          />
        </div>

        {saveMessage && (
          <div className={`rounded-md px-4 py-3 text-sm ${saveMessage.type === 'success' ? 'bg-green-50 text-green-800 border border-green-200' : 'bg-red-50 text-red-800 border border-red-200'}`}>
            {saveMessage.text}
          </div>
        )}

        <div className="flex justify-end pt-4">
          <Button onClick={handleSave} disabled={onUpdate.isPending}>
            {onUpdate.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            <Save className="mr-2 h-4 w-4" />
            Save Changes
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}

// ── Attributes Tab ──
function AttributesTab({ productId, attributeGroupId, isConfigurable }: { productId: string; attributeGroupId: string; isConfigurable: boolean }) {
  const { data: groupAttributes, isLoading: loadingGroupAttrs } = useGroupAttributes(attributeGroupId);
  const { data: productAttributes, isLoading: loadingProdAttrs } = useAdminProductAttributes(productId);
  const { data: productOptions } = useAdminProductOptions(productId);
  const setAttributes = useAdminSetProductAttributes();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [values, setValues] = useState<Record<string, string>>({});
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [multiValues, setMultiValues] = useState<Record<string, string[]>>({});
  const [initialized, setInitialized] = useState(false);
  const [saveMessage, setSaveMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  // Initialize from existing product attributes
  if (!initialized && productAttributes && !loadingProdAttrs) {
    const v: Record<string, string> = {};
    const mv: Record<string, string[]> = {};
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    productAttributes.forEach((attr: any) => {
      if (attr.values && attr.values.length > 0) {
        mv[attr.attribute_id] = attr.values;
      } else {
        v[attr.attribute_id] = attr.value || '';
      }
    });
    setValues(v);
    setMultiValues(mv);
    setInitialized(true);
  }

  if (!attributeGroupId) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Product Attributes</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <Tags className="h-12 w-12 text-muted-foreground/30" />
            <p className="mt-4 text-sm text-muted-foreground">
              No attribute group selected. Go to the General tab and select an attribute group first.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (loadingGroupAttrs || loadingProdAttrs) {
    return <Skeleton className="h-64 w-full" />;
  }

  const allAttrs: Attribute[] = groupAttributes || [];

  // For configurable products, hide variation-eligible attributes (select/multi_select/color with options).
  // These are managed on the Variations tab, not as product-wide specs.
  const isVariationEligible = (attr: Attribute) =>
    (attr.type === 'select' || attr.type === 'multi_select' || attr.type === 'color') &&
    attr.options &&
    attr.options.length > 0;

  const attrs = isConfigurable
    ? allAttrs.filter((attr) => !isVariationEligible(attr))
    : allAttrs;
  const hiddenAttrs = isConfigurable
    ? allAttrs.filter((attr) => isVariationEligible(attr))
    : [];
  const hiddenCount = hiddenAttrs.length;

  const handleSave = () => {
    setSaveMessage(null);
    const attrPayload = attrs.map((attr) => ({
      attribute_id: attr.id,
      value: values[attr.id] || '',
      values: multiValues[attr.id] || undefined,
    }));
    setAttributes.mutate(
      { productId, attributes: attrPayload },
      {
        onSuccess: () => setSaveMessage({ type: 'success', text: 'Attributes saved successfully.' }),
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        onError: (err: any) => setSaveMessage({ type: 'error', text: err?.response?.data?.error || err?.message || 'Failed to save attributes.' }),
      }
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Product Attributes</CardTitle>
        <CardDescription>
          {isConfigurable
            ? 'Product-wide specs that apply to all variants (e.g., Material, Brand)'
            : 'Fill in attribute values defined by the selected attribute group'}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {hiddenCount > 0 && (
          <div className="rounded-md bg-muted/50 px-4 py-3 text-sm text-muted-foreground">
            {hiddenCount} attribute{hiddenCount > 1 ? 's' : ''} ({hiddenAttrs.map(a => a.name.toLowerCase()).join(', ')}) hidden — managed as variation options on the Variations tab.
          </div>
        )}
        {attrs.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {hiddenCount > 0
              ? 'All attributes in this group are used for variations. Add non-variation attributes (e.g., Material, Brand) to the attribute group to manage them here.'
              : 'No attributes defined in this group. Add attributes to the group first.'}
          </p>
        ) : (
          attrs.map((attr) => (
            <div key={attr.id} className="space-y-2">
              <Label>
                {attr.name}
                {attr.required && <span className="ml-1 text-destructive">*</span>}
                {attr.unit && (
                  <span className="ml-1 text-xs text-muted-foreground">({attr.unit})</span>
                )}
              </Label>
              {attr.type === 'select' && attr.options ? (
                <select
                  value={values[attr.id] || ''}
                  onChange={(e) => setValues({ ...values, [attr.id]: e.target.value })}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="">Select...</option>
                  {attr.options.map((opt) => (
                    <option key={opt} value={opt}>
                      {opt}
                    </option>
                  ))}
                </select>
              ) : attr.type === 'multi_select' && attr.options ? (
                <div className="space-y-2">
                  <div className="flex flex-wrap gap-1.5">
                    {(multiValues[attr.id] || []).map((val) => (
                      <Badge key={val} variant="secondary" className="gap-1 pr-1">
                        {val}
                        <button
                          type="button"
                          onClick={() =>
                            setMultiValues({
                              ...multiValues,
                              [attr.id]: (multiValues[attr.id] || []).filter((v) => v !== val),
                            })
                          }
                          className="ml-1 rounded-full p-0.5 hover:bg-muted"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </Badge>
                    ))}
                  </div>
                  <select
                    value=""
                    onChange={(e) => {
                      const val = e.target.value;
                      if (val && !(multiValues[attr.id] || []).includes(val)) {
                        setMultiValues({
                          ...multiValues,
                          [attr.id]: [...(multiValues[attr.id] || []), val],
                        });
                      }
                    }}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  >
                    <option value="">Add a value...</option>
                    {attr.options
                      .filter((opt) => !(multiValues[attr.id] || []).includes(opt))
                      .map((opt) => (
                        <option key={opt} value={opt}>
                          {opt}
                        </option>
                      ))}
                  </select>
                </div>
              ) : attr.type === 'bool' ? (
                <select
                  value={values[attr.id] || ''}
                  onChange={(e) => setValues({ ...values, [attr.id]: e.target.value })}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="">Select...</option>
                  <option value="true">Yes</option>
                  <option value="false">No</option>
                </select>
              ) : attr.type === 'color' ? (
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    value={values[attr.id] || '#000000'}
                    onChange={(e) => setValues({ ...values, [attr.id]: e.target.value })}
                    className="h-10 w-10 cursor-pointer rounded border"
                  />
                  <Input
                    value={values[attr.id] || ''}
                    onChange={(e) => setValues({ ...values, [attr.id]: e.target.value })}
                    placeholder="#000000"
                    className="flex-1"
                  />
                </div>
              ) : (
                <Input
                  type={attr.type === 'number' ? 'number' : 'text'}
                  value={values[attr.id] || ''}
                  onChange={(e) => setValues({ ...values, [attr.id]: e.target.value })}
                  placeholder={`Enter ${attr.name.toLowerCase()}...`}
                />
              )}
            </div>
          ))
        )}

        {attrs.length > 0 && (
          <>
            {saveMessage && (
              <div className={`rounded-md px-4 py-3 text-sm ${saveMessage.type === 'success' ? 'bg-green-50 text-green-800 border border-green-200' : 'bg-red-50 text-red-800 border border-red-200'}`}>
                {saveMessage.text}
              </div>
            )}
            <div className="flex justify-end pt-4">
              <Button onClick={handleSave} disabled={setAttributes.isPending}>
                {setAttributes.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                <Save className="mr-2 h-4 w-4" />
                Save Attributes
              </Button>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

// ── Variations Tab (Configurable products only) ──
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function VariationsTab({ productId, product, attributeGroupId }: { productId: string; product: any; attributeGroupId: string }) {
  const { data: options, isLoading: loadingOptions } = useAdminProductOptions(productId);
  const { data: variants, isLoading: loadingVariants } = useAdminProductVariants(productId);
  const { data: groupAttributes } = useGroupAttributes(attributeGroupId);
  const addOption = useAdminAddOption();
  const removeOption = useAdminRemoveOption();
  const generateVariants = useAdminGenerateVariants();
  const updateVariant = useAdminUpdateVariant();

  // New option form state
  const [showAddOption, setShowAddOption] = useState(false);
  const [newOptionName, setNewOptionName] = useState('');
  const [newOptionValues, setNewOptionValues] = useState('');

  // Confirm dialog
  const [confirmGenerate, setConfirmGenerate] = useState(false);

  // Attribute-to-option selection state: tracks which values are selected per attribute
  const [selectedAttrValues, setSelectedAttrValues] = useState<Record<string, string[]>>({});

  // Filter to variation-eligible attributes (select, multi_select, color have predefined options)
  const variationAttributes = (groupAttributes || []).filter(
    (attr: Attribute) =>
      (attr.type === 'select' || attr.type === 'multi_select' || attr.type === 'color') &&
      attr.options &&
      attr.options.length > 0
  );

  // Check if an option with the given name already exists
  const isOptionAlreadyAdded = (attrName: string) =>
    (options || []).some(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (opt: any) => opt.name.toLowerCase() === attrName.toLowerCase()
    );

  // Toggle a single value in the selection for an attribute
  const toggleAttrValue = (attrName: string, value: string, allOptions: string[]) => {
    setSelectedAttrValues((prev) => {
      const current = prev[attrName] ?? allOptions;
      const next = current.includes(value)
        ? current.filter((v) => v !== value)
        : [...current, value];
      return { ...prev, [attrName]: next };
    });
  };

  // Add attribute as a product option
  const handleAddAttrAsOption = (attr: Attribute) => {
    const selected = selectedAttrValues[attr.name] ?? attr.options ?? [];
    if (selected.length === 0) return;

    addOption.mutate({
      productId,
      data: {
        name: attr.name,
        values: selected.map((v, i) => ({ value: v, sort_order: i })),
      },
    });
  };

  const handleAddOption = () => {
    if (!newOptionName.trim()) return;
    const vals = newOptionValues
      .split(',')
      .map((v) => v.trim())
      .filter(Boolean)
      .map((v, i) => ({ value: v, sort_order: i }));
    if (vals.length === 0) return;

    addOption.mutate(
      { productId, data: { name: newOptionName.trim(), values: vals } },
      {
        onSuccess: () => {
          setNewOptionName('');
          setNewOptionValues('');
          setShowAddOption(false);
        },
      }
    );
  };

  const handleGenerateVariants = () => {
    generateVariants.mutate(productId);
    setConfirmGenerate(false);
  };

  return (
    <div className="space-y-6">
      {/* Options Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Options</CardTitle>
              <CardDescription>
                Define product options like Size, Color, Material
              </CardDescription>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowAddOption(!showAddOption)}
            >
              <Plus className="mr-1 h-4 w-4" />
              Custom Option
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* From Attribute Group section */}
          {variationAttributes.length > 0 && (
            <div className="space-y-3">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                From Attribute Group
              </p>
              {variationAttributes.map((attr: Attribute) => {
                const alreadyAdded = isOptionAlreadyAdded(attr.name);
                const attrOptions = attr.options || [];
                const selected = selectedAttrValues[attr.name] ?? attrOptions;

                return (
                  <div
                    key={attr.id}
                    className={`rounded-lg border p-3 space-y-2 ${alreadyAdded ? 'bg-muted/30' : ''}`}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <p className="font-medium text-sm">{attr.name}</p>
                        <Badge variant="outline" className="text-xs">{attr.type}</Badge>
                      </div>
                      {alreadyAdded ? (
                        <Badge variant="secondary" className="text-xs gap-1">
                          <Check className="h-3 w-3" />
                          Added
                        </Badge>
                      ) : (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleAddAttrAsOption(attr)}
                          disabled={addOption.isPending || selected.length === 0}
                        >
                          {addOption.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
                          Add to Options
                        </Button>
                      )}
                    </div>
                    {!alreadyAdded && (
                      <div className="flex flex-wrap gap-1.5 max-h-32 overflow-y-auto">
                        {attrOptions.map((val) => {
                          const isSelected = selected.includes(val);
                          return (
                            <button
                              key={val}
                              type="button"
                              onClick={() => toggleAttrValue(attr.name, val, attrOptions)}
                              className={`inline-flex items-center rounded-md border px-2.5 py-1 text-xs font-medium transition-colors ${
                                isSelected
                                  ? 'border-primary bg-primary/10 text-primary'
                                  : 'border-input bg-background text-muted-foreground hover:bg-muted'
                              }`}
                            >
                              {isSelected && <Check className="mr-1 h-3 w-3" />}
                              {val}
                            </button>
                          );
                        })}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}

          {/* Manual Add Option form */}
          {showAddOption && (
            <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                Custom Option
              </p>
              <div className="space-y-2">
                <Label>Option Name</Label>
                <Input
                  value={newOptionName}
                  onChange={(e) => setNewOptionName(e.target.value)}
                  placeholder="e.g. Size, Color, Material"
                />
              </div>
              <div className="space-y-2">
                <Label>Values (comma-separated)</Label>
                <Input
                  value={newOptionValues}
                  onChange={(e) => setNewOptionValues(e.target.value)}
                  placeholder="e.g. Small, Medium, Large"
                />
              </div>
              <div className="flex gap-2 justify-end">
                <Button variant="ghost" size="sm" onClick={() => setShowAddOption(false)}>
                  Cancel
                </Button>
                <Button size="sm" onClick={handleAddOption} disabled={addOption.isPending}>
                  {addOption.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
                  Add Option
                </Button>
              </div>
            </div>
          )}

          {/* Existing options list */}
          {loadingOptions ? (
            <Skeleton className="h-20 w-full" />
          ) : !options || options.length === 0 ? (
            variationAttributes.length === 0 && (
              <p className="text-sm text-muted-foreground py-4 text-center">
                No options defined yet. Add options to configure product variations.
              </p>
            )
          ) : (
            <div className="space-y-3">
              {options.length > 0 && (
                <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Active Options
                </p>
              )}
              {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
              {options.map((opt: any) => {
                const isFromAttr = variationAttributes.some(
                  (attr: Attribute) => attr.name.toLowerCase() === opt.name.toLowerCase()
                );
                return (
                  <div
                    key={opt.id}
                    className="flex items-center justify-between rounded-lg border p-3"
                  >
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="font-medium text-sm">{opt.name}</p>
                        {isFromAttr && (
                          <Badge variant="outline" className="text-[10px] text-muted-foreground">
                            from attributes
                          </Badge>
                        )}
                      </div>
                      <div className="flex flex-wrap gap-1 mt-1">
                        {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                        {(opt.values || []).map((v: any) => (
                          <Badge key={v.id} variant="secondary" className="text-xs">
                            {v.value}
                          </Badge>
                        ))}
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-destructive"
                      onClick={() =>
                        removeOption.mutate({ productId, optionId: opt.id })
                      }
                      disabled={removeOption.isPending}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                );
              })}
            </div>
          )}

          {options && options.length > 0 && (
            <div className="flex justify-end pt-2">
              <Button onClick={() => setConfirmGenerate(true)}>
                <Wand2 className="mr-2 h-4 w-4" />
                Generate Variants
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Variants Card */}
      <Card>
        <CardHeader>
          <CardTitle>Variants</CardTitle>
          <CardDescription>
            Manage pricing, stock, and details for each variant
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loadingVariants ? (
            <Skeleton className="h-32 w-full" />
          ) : !variants || variants.length === 0 ? (
            <p className="text-sm text-muted-foreground py-8 text-center">
              No variants generated yet. Add options and click "Generate Variants" to create them.
            </p>
          ) : (
            <div className="space-y-3">
              {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
              {variants.map((variant: any) => (
                <VariantRow
                  key={variant.id}
                  variant={variant}
                  productId={productId}
                  onUpdate={updateVariant}
                />
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Generate Confirm Dialog */}
      <Dialog open={confirmGenerate} onOpenChange={setConfirmGenerate}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Generate Variants</DialogTitle>
            <DialogDescription>
              This will generate new variants from the cartesian product of all options.
              Existing variants with the same SKU will be kept. New variants will start with
              the base product price and 0 stock.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setConfirmGenerate(false)}>
              Cancel
            </Button>
            <Button onClick={handleGenerateVariants} disabled={generateVariants.isPending}>
              {generateVariants.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Generate
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

// ── Variant Row (inline editing) ──
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function VariantRow({ variant, productId, onUpdate }: { variant: any; productId: string; onUpdate: any }) {
  const [editing, setEditing] = useState(false);
  const [price, setPrice] = useState(variant.price_cents);
  const [stock, setStock] = useState(variant.stock);
  const [sku, setSku] = useState(variant.sku);
  const [showImages, setShowImages] = useState(false);
  const [variantImages, setVariantImages] = useState<string[]>(variant.image_urls || []);

  const handleSave = () => {
    onUpdate.mutate(
      {
        productId,
        variantId: variant.id,
        data: { price_cents: price, stock, name: variant.name, image_urls: variantImages },
      },
      { onSuccess: () => setEditing(false) }
    );
  };

  const handleImageChange = (urls: string[]) => {
    setVariantImages(urls);
    // Auto-save variant images
    onUpdate.mutate({
      productId,
      variantId: variant.id,
      data: { name: variant.name, image_urls: urls },
    });
  };

  return (
    <div className="rounded-lg border p-3 space-y-3">
      {/* Row 1: Variant name + option badges */}
      <div className="flex items-start justify-between gap-2">
        <div>
          <p className="font-medium text-sm">{variant.name}</p>
          <div className="flex flex-wrap gap-1 mt-0.5">
            {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
            {(variant.option_values || []).map((ov: any) => (
              <Badge key={ov.option_value_id} variant="outline" className="text-xs">
                {ov.option_name}: {ov.value}
              </Badge>
            ))}
          </div>
        </div>
        <div className="flex items-center gap-1.5 shrink-0">
          <Button variant="ghost" size="sm" className="h-7 px-2" onClick={() => setShowImages(!showImages)}>
            <ImageIcon className="h-3.5 w-3.5" />
            {variantImages.length > 0 && <span className="ml-1 text-xs">{variantImages.length}</span>}
          </Button>
          <Button variant="ghost" size="sm" className="h-7 px-2" onClick={() => setEditing(!editing)}>
            {editing ? 'Cancel' : 'Edit'}
          </Button>
        </div>
      </div>

      {/* Row 2: Price, stock, status display */}
      {!editing ? (
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm font-mono">${(variant.price_cents / 100).toFixed(2)}</span>
          <Badge variant={variant.stock > 0 ? 'secondary' : 'destructive'} className="text-xs">
            Stock: {variant.stock}
          </Badge>
          <Badge variant={variant.is_active ? 'default' : 'outline'} className="text-xs">
            {variant.is_active ? 'Active' : 'Inactive'}
          </Badge>
        </div>
      ) : (
        <div className="space-y-2 pt-1 border-t">
          <div className="grid grid-cols-2 gap-2">
            <div className="space-y-1">
              <Label className="text-xs text-muted-foreground">Price (cents)</Label>
              <Input
                type="number"
                value={price}
                onChange={(e) => setPrice(Number(e.target.value))}
                className="h-8 text-sm"
              />
            </div>
            <div className="space-y-1">
              <Label className="text-xs text-muted-foreground">Stock</Label>
              <Input
                type="number"
                value={stock}
                onChange={(e) => setStock(Number(e.target.value))}
                className="h-8 text-sm"
              />
            </div>
          </div>
          <div className="flex justify-end gap-2">
            <Button size="sm" onClick={handleSave} disabled={onUpdate.isPending}>
              {onUpdate.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
              Save
            </Button>
          </div>
        </div>
      )}

      {/* Collapsible image uploader */}
      {showImages && (
        <div className="pt-2 border-t">
          <ImageUploader
            value={variantImages}
            onChange={handleImageChange}
            maxFiles={3}
            maxSizeMB={5}
            ownerType="variant"
          />
        </div>
      )}
    </div>
  );
}
