<script lang="ts">
    import * as Sidebar from "#lib/components/ui/sidebar/index.js";
    import logo from "#lib/assets/yokify.svg"
</script>

<Sidebar.Root>
    <Sidebar.Header>
        <Sidebar.GroupLabel class="py-6 flex justify-start">
            <img class="w-12 h-12" src={logo} alt="Yokify Logo" />
            <h1 class="text-3xl font-bold">Yōkify</h1>
        </Sidebar.GroupLabel>
    </Sidebar.Header>
    <Sidebar.Separator />
    <Sidebar.Content>
        <Sidebar.Menu>
            <Sidebar.MenuButton class="text-center text-xl font-light flex px-6">
                Functions
            </Sidebar.MenuButton>
        </Sidebar.Menu>
    </Sidebar.Content>
    <Sidebar.Footer />
</Sidebar.Root>